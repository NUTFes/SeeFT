import os
import socketserver
from http.server import SimpleHTTPRequestHandler

PORT = 45029
WEB_ROOT = "./build/web"


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
        return super().send_head()


with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()
