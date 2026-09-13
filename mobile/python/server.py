import datetime
import email.utils
import gzip
import io
import os
import threading
from http import HTTPStatus
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer

PORT = 45029
WEB_ROOT = "./build/web"

# gzip で圧縮して返す拡張子。本番は Cloudflare がスマホまでの区間を圧縮するが、
# ここから Cloudflare までは生のサイズで流れるため、main.dart.js や同梱フォントはここで圧縮して送る。
# 画像は元から圧縮されているので対象にしない
COMPRESSIBLE_EXTENSIONS = {".html", ".js", ".json", ".css", ".wasm", ".ttf", ".otf", ".svg", ".txt"}

_gzip_cache = {}
_gzip_lock = threading.Lock()


def gzipped_body(path, stat):
    # ビルド後にファイルは変わらないので、ファイルごとに 1 回だけ圧縮してメモリに持つ
    key = (path, stat.st_mtime, stat.st_size)
    body = _gzip_cache.get(key)
    if body is None:
        with _gzip_lock:
            body = _gzip_cache.get(key)
            if body is None:
                with open(path, "rb") as f:
                    body = gzip.compress(f.read(), compresslevel=6)
                _gzip_cache[key] = body
    return body


class MyHandler(SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=WEB_ROOT, **kwargs)

    # アプリは URL に # を付けない形式（PathUrlStrategy）なので、ログイン後は /layout のような URL になる。
    # その状態で再読み込みされると該当するファイルが無く 404 になるため、
    # ファイルが無く拡張子も無いパスには index.html を返し、どの画面を出すかはアプリに任せる。
    # 拡張子の付いたパス（.js など）は、ファイルが欠けていれば今まで通り 404 にして気づけるようにする
    def send_head(self):
        path = self.translate_path(self.path)
        if not os.path.exists(path) and os.path.splitext(path)[1] == "":
            self.path = "/index.html"
            path = self.translate_path(self.path)

        self.compressible = os.path.splitext(path)[1] in COMPRESSIBLE_EXTENSIONS
        if self.compressible and os.path.isfile(path) and self.accepts_gzip():
            return self.send_gzip_head(path)
        return super().send_head()

    def end_headers(self):
        # 圧縮できるファイルは、Accept-Encoding によって返す中身が変わることをキャッシュに伝える
        if getattr(self, "compressible", False):
            self.send_header("Vary", "Accept-Encoding")
        super().end_headers()

    def accepts_gzip(self):
        # gzip に付いた q の値が 0 なら（q=0.0 や Q=0 も含む）、圧縮を受け付けない指定として扱う
        for part in self.headers.get("Accept-Encoding", "").split(","):
            coding, *params = part.split(";")
            if coding.strip().lower() != "gzip":
                continue
            quality = 1.0
            for param in params:
                name, _, value = param.partition("=")
                if name.strip().lower() == "q":
                    try:
                        quality = float(value)
                    except ValueError:
                        return False
            return 0 < quality <= 1
        return False

    def send_gzip_head(self, path):
        stat = os.stat(path)
        if self.not_modified(stat):
            self.send_response(HTTPStatus.NOT_MODIFIED)
            self.end_headers()
            return None
        body = gzipped_body(path, stat)
        self.send_response(HTTPStatus.OK)
        self.send_header("Content-Type", self.guess_type(path))
        self.send_header("Content-Encoding", "gzip")
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Last-Modified", self.date_time_string(stat.st_mtime))
        self.end_headers()
        return io.BytesIO(body)

    def not_modified(self, stat):
        # SimpleHTTPRequestHandler.send_head と同じ判定で、更新が無ければ 304 を返す
        if "If-Modified-Since" not in self.headers or "If-None-Match" in self.headers:
            return False
        try:
            since = email.utils.parsedate_to_datetime(self.headers["If-Modified-Since"])
        except (TypeError, IndexError, OverflowError, ValueError):
            return False
        if since.tzinfo is None:
            since = since.replace(tzinfo=datetime.timezone.utc)
        if since.tzinfo is not datetime.timezone.utc:
            return False
        modified = datetime.datetime.fromtimestamp(stat.st_mtime, datetime.timezone.utc)
        return modified.replace(microsecond=0) <= since


# ThreadingHTTPServer はリクエストごとにスレッドを立てる。
# 1 本のスレッドで順番に返すと、初めて開く人が重なったときに後ろの人が待たされるため
with ThreadingHTTPServer(("", PORT), MyHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()
