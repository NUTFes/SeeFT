#!/usr/bin/env python3
"""
manifest.json + 作業写真 → セクションmp4 動画

使い方:
  python3 scripts/assemble_video.py docs/manuals/01_44th_配線マニュアル --section 0
  python3 scripts/assemble_video.py docs/manuals/01_44th_配線マニュアル

画面設計: 作業写真主役型
  - 画像を中央に最大表示（アスペクト比維持、余白黒）
  - 字幕を下部に焼き込み
  - キャラ立ち絵なし
"""

import argparse
import json
import os
import subprocess
from collections import defaultdict

from PIL import Image, ImageDraw, ImageFont

VIDEO_W = 1280
VIDEO_H = 720
FRAMERATE = 25

# 日本語フォント（drawtext 用）
JP_FONT_SRC = "/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc"

# 字幕スタイル
SUB_FONT_SIZE = 30
SUB_MAX_CHARS = 26
SUB_COLOR_RGB = {
    "Expert": (155, 231, 212),
    "Novice": (255, 214, 153),
}
SUB_BG_RGBA = (0, 0, 0, 200)
def ts_srt(s: float) -> str:
    h = int(s // 3600)
    m = int((s % 3600) // 60)
    sec = s - h * 3600 - m * 60
    return f"{h:02d}:{m:02d}:{int(sec):02d},{int((sec - int(sec)) * 1000):03d}"


def ts_vtt(s: float) -> str:
    h = int(s // 3600)
    m = int((s % 3600) // 60)
    sec = s - h * 3600 - m * 60
    return f"{h:02d}:{m:02d}:{int(sec):02d}.{int((sec - int(sec)) * 1000):03d}"


def build_srt(entries: list, srt_path: str):
    with open(srt_path, "w", encoding="utf-8") as f:
        t = 0.0
        for i, e in enumerate(entries, start=1):
            start, end = t, t + e["duration_sec"]
            label = f"[{e['speaker']}] {e['text']}"
            f.write(f"{i}\n{ts_srt(start)} --> {ts_srt(end)}\n{label}\n\n")
            t = end


def wrap_ja(text: str, max_chars: int = SUB_MAX_CHARS) -> str:
    """日本語テキストを max_chars ごとに強制改行（句読点優先）"""
    out_lines = []
    buf = ""
    for ch in text:
        buf += ch
        if ch in "。！？":
            out_lines.append(buf)
            buf = ""
        elif len(buf) >= max_chars:
            # 近傍の句読点で切れるなら切る
            cut = len(buf)
            for i in range(len(buf) - 1, max(0, len(buf) - 6), -1):
                if buf[i] in "、。，．":
                    cut = i + 1
                    break
            out_lines.append(buf[:cut])
            buf = buf[cut:]
    if buf:
        out_lines.append(buf)
    return "\n".join(out_lines)


def render_subtitle_png(text: str, color_rgb: tuple, png_path: str, w: int = VIDEO_W, h: int = VIDEO_H):
    """字幕PNGを生成 (透明背景、下部に黒半透明ボックス+テキスト)"""
    img = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    font = ImageFont.truetype(JP_FONT_SRC, size=SUB_FONT_SIZE)

    lines = text.split("\n")
    line_h = SUB_FONT_SIZE + 12
    total_h = len(lines) * line_h

    # 下部中央配置
    pad_x, pad_y = 20, 12
    base_y = h - total_h - 50

    for i, line in enumerate(lines):
        bbox = draw.textbbox((0, 0), line, font=font)
        tw = bbox[2] - bbox[0]
        th = bbox[3] - bbox[1]
        x = (w - tw) // 2
        y = base_y + i * line_h

        # 半透明ボックス
        draw.rounded_rectangle(
            [x - pad_x, y - pad_y, x + tw + pad_x, y + line_h - pad_y],
            radius=8, fill=SUB_BG_RGBA,
        )
        # テキスト
        draw.text((x, y - bbox[1]), line, fill=(*color_rgb, 255), font=font)

    img.save(png_path, "PNG")


def build_subtitle_pngs(entries: list, work_dir: str) -> list:
    """各ターンの字幕PNGを生成し、パスと再生時間のリストを返す"""
    sub_dir = os.path.join(work_dir, "subs")
    os.makedirs(sub_dir, exist_ok=True)
    result = []
    t = 0.0
    for idx, e in enumerate(entries):
        start, end = t, t + e["duration_sec"]
        wrapped = wrap_ja(f"[{e['speaker']}] {e['text']}")
        png_path = os.path.join(sub_dir, f"sub_{idx:03d}.png")
        color = SUB_COLOR_RGB.get(e["speaker"], (255, 255, 255))
        render_subtitle_png(wrapped, color, png_path)
        result.append((png_path, start, end))
        t = end
    return result


def build_vtt(entries: list, vtt_path: str):
    with open(vtt_path, "w", encoding="utf-8") as f:
        f.write("WEBVTT\n\n")
        t = 0.0
        for i, e in enumerate(entries, start=1):
            start, end = t, t + e["duration_sec"]
            cue_class = "expert" if e["speaker"] == "Expert" else "novice"
            f.write(f"{i}\n{ts_vtt(start)} --> {ts_vtt(end)}\n")
            f.write(f"<c.{cue_class}>[{e['speaker']}] {e['text']}</c>\n\n")
            t = end


def concat_audio(entries: list, manual_dir: str, work_dir: str, out_wav: str):
    list_path = os.path.join(work_dir, "audio_list.txt")
    with open(list_path, "w") as f:
        for e in entries:
            abs_path = os.path.abspath(os.path.join(manual_dir, e["audio_path"]))
            f.write(f"file '{abs_path}'\n")
    subprocess.run([
        "ffmpeg", "-y", "-loglevel", "error",
        "-f", "concat", "-safe", "0", "-i", list_path,
        "-c:a", "pcm_s16le", "-ar", "24000", "-ac", "1",
        out_wav,
    ], check=True)


def fill_image_timeline(entries: list):
    current = None
    timeline = []
    for e in entries:
        if e.get("image"):
            current = e["image"]
        timeline.append((current, e["duration_sec"]))
    return timeline


def merge_consecutive(timeline: list):
    merged = []
    for img, dur in timeline:
        if merged and merged[-1][0] == img:
            merged[-1][1] += dur
        else:
            merged.append([img, dur])
    return merged


def ensure_black_png(work_dir: str) -> str:
    """画像なしセクション用の黒背景PNGを用意 (1度だけ生成)"""
    path = os.path.join(work_dir, "black.png")
    if not os.path.exists(path):
        Image.new("RGB", (VIDEO_W, VIDEO_H), (0, 0, 0)).save(path)
    return path


def build_section_video(entries: list, manual_dir: str, out_path: str, work_dir: str):
    os.makedirs(work_dir, exist_ok=True)
    images_dir = os.path.join(manual_dir, "images")

    # 1. 音声トラック
    combined_audio = os.path.join(work_dir, "audio.wav")
    concat_audio(entries, manual_dir, work_dir, combined_audio)

    # 2. サイドカー字幕 (SRT + VTT)
    out_dir = os.path.dirname(out_path)
    base = os.path.splitext(os.path.basename(out_path))[0]
    build_srt(entries, os.path.join(out_dir, f"{base}.srt"))
    build_vtt(entries, os.path.join(out_dir, f"{base}.vtt"))

    # 3. 画像タイムライン (連続する同一画像を統合)
    merged = merge_consecutive(fill_image_timeline(entries))

    # 4. 字幕PNG (ターンごと、透明背景+テキストボックス)
    sub_specs = build_subtitle_pngs(entries, work_dir)

    # 5. 単一の filter_complex で: 画像concat → setpts正規化 → 字幕overlay連鎖
    black_png = ensure_black_png(work_dir)

    cmd = ["ffmpeg", "-y", "-loglevel", "error"]

    # 画像入力群 (各画像を -loop 1 -t で指定秒数再生)
    image_input_count = len(merged)
    for img, dur in merged:
        src = os.path.join(images_dir, img) if img and os.path.exists(os.path.join(images_dir, img)) else black_png
        cmd.extend(["-loop", "1", "-t", f"{dur:.3f}", "-framerate", str(FRAMERATE), "-i", src])

    # 字幕PNG入力群
    for png_path, _, _ in sub_specs:
        cmd.extend(["-i", png_path])

    # 音声入力
    cmd.extend(["-i", combined_audio])

    sub_input_start = image_input_count
    audio_input_idx = sub_input_start + len(sub_specs)

    # filter_complex 構築
    filter_parts = []

    # 各画像を固定サイズにスケール + アスペクト維持 + FPS統一
    for i in range(image_input_count):
        filter_parts.append(
            f"[{i}:v]scale={VIDEO_W}:{VIDEO_H}:force_original_aspect_ratio=decrease,"
            f"pad={VIDEO_W}:{VIDEO_H}:(ow-iw)/2:(oh-ih)/2:color=black,setsar=1,"
            f"fps={FRAMERATE},format=yuv420p[vimg{i}]"
        )

    # concat フィルタで画像セグメントを連結 (PTS は自動的に連続化される)
    concat_inputs = "".join(f"[vimg{i}]" for i in range(image_input_count))
    filter_parts.append(f"{concat_inputs}concat=n={image_input_count}:v=1:a=0[vconcat]")

    # 0 起点に正規化
    filter_parts.append("[vconcat]setpts=PTS-STARTPTS[vbase]")

    # 字幕 overlay 連鎖 (各ターンの PNG を時間窓で重ねる)
    last_label = "vbase"
    for idx, (_, start, end) in enumerate(sub_specs):
        input_idx = sub_input_start + idx
        next_label = f"vov{idx + 1}"
        filter_parts.append(
            f"[{last_label}][{input_idx}:v]"
            f"overlay=enable='between(t\\,{start:.3f}\\,{end:.3f})'"
            f"[{next_label}]"
        )
        last_label = next_label

    filter_complex = ";".join(filter_parts)

    cmd.extend([
        "-filter_complex", filter_complex,
        "-map", f"[{last_label}]",
        "-map", f"{audio_input_idx}:a",
        "-c:v", "libx264", "-preset", "medium", "-crf", "23",
        "-c:a", "aac", "-b:a", "192k",
        "-pix_fmt", "yuv420p",
        "-shortest",
        out_path,
    ])
    subprocess.run(cmd, check=True)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("manual_dir")
    ap.add_argument("--section", type=int, help="対象セクションindex (0-based)")
    args = ap.parse_args()

    manual_dir = args.manual_dir.rstrip("/")
    manifest_path = os.path.join(manual_dir, "audio", "manifest.json")
    videos_dir = os.path.join(manual_dir, "videos")
    os.makedirs(videos_dir, exist_ok=True)

    with open(manifest_path, encoding="utf-8") as f:
        manifest = json.load(f)

    by_section = defaultdict(list)
    for e in manifest["entries"]:
        by_section[e["section_idx"]].append(e)

    targets = [args.section] if args.section is not None else sorted(by_section.keys())

    for i in targets:
        if i not in by_section:
            print(f"section {i} not in manifest")
            continue
        print(f"=== section {i}: {len(by_section[i])} turns ===")
        work_dir = os.path.join(videos_dir, f"_work_sec{i:02d}")
        out_path = os.path.join(videos_dir, f"section_{i:02d}.mp4")
        build_section_video(by_section[i], manual_dir, out_path, work_dir)
        size_kb = os.path.getsize(out_path) // 1024
        print(f"  -> {out_path} ({size_kb}KB)")


if __name__ == "__main__":
    main()
