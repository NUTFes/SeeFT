#!/usr/bin/env python3
"""
slide_api.html に解説動画セクションを追加した manual_with_videos.html を生成。

使い方:
  python3 scripts/build_manual_with_videos.py docs/manuals/01_44th_配線マニュアル

依存する事前生成物:
  - slide_api.html (generate_manual_slide.py の出力)
  - dialogue.json (generate_dialogue_script.py の出力)
  - videos/section_NN.mp4 (assemble_video.py の出力、1本以上)

slide_api.html そのものは触らないので、AI による再生成で上書きされても
このスクリプトを再実行すれば最新のスライド + 動画で再構築できる。
"""

import json
import os
import re
import sys


def build_video_section(entries: list) -> str:
    """動画セクションのHTMLを組み立てる"""
    if not entries:
        return ""

    buttons = "\n    ".join(
        f'<button class="chapter-btn" data-video="videos/{e["file"]}" '
        f'data-index="{e["index"]}">'
        f'{e["index"] + 1}. {e["title"]}'
        f'</button>'
        for e in entries
    )

    default_video = f"videos/{entries[0]['file']}"

    return f"""
<section id="video-gallery" style="background:var(--sub-light-1,#f4fbf8);padding:clamp(1rem,3vw,2rem);border-bottom:4px solid var(--main,#009688);scroll-snap-align:start;">
  <h2 style="color:var(--main,#009688);font-size:clamp(18px,3vw,24px);margin:0 0 12px;display:flex;align-items:center;gap:8px;">
    <span aria-hidden="true">🎬</span> 解説動画
  </h2>
  <div style="background:#000;border-radius:8px;overflow:hidden;max-width:900px;margin:0 auto 16px;">
    <video id="manual-video" controls preload="metadata" style="width:100%;aspect-ratio:16/9;display:block;">
      <source src="{default_video}" type="video/mp4">
    </video>
  </div>
  <div class="chapter-list" role="tablist" aria-label="セクション選択" style="display:flex;flex-wrap:wrap;gap:8px;justify-content:center;">
    {buttons}
  </div>
  <style>
    .chapter-btn {{
      padding: 8px 14px;
      border: 2px solid var(--main, #009688);
      background: white;
      color: var(--main, #009688);
      border-radius: 8px;
      cursor: pointer;
      font-size: clamp(12px, 2vw, 14px);
      font-family: inherit;
      font-weight: 600;
      transition: background 0.15s, color 0.15s;
    }}
    .chapter-btn:hover {{ background: var(--sub-light-2, #cce8e2); }}
    .chapter-btn.active {{ background: var(--main, #009688); color: white; }}
    .chapter-btn:focus-visible {{ outline: 3px solid #F3AE56; outline-offset: 2px; }}
  </style>
  <script>
    (function() {{
      var video = document.getElementById('manual-video');
      var source = video.querySelector('source');
      var buttons = document.querySelectorAll('.chapter-btn');
      if (buttons.length > 0) buttons[0].classList.add('active');
      buttons.forEach(function(btn) {{
        btn.addEventListener('click', function() {{
          buttons.forEach(function(b) {{ b.classList.remove('active'); }});
          btn.classList.add('active');
          source.src = btn.dataset.video;
          video.load();
          video.play().catch(function(){{}});
        }});
      }});
    }})();
  </script>
</section>
"""


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/build_manual_with_videos.py <manual_dir>")
        sys.exit(1)

    manual_dir = sys.argv[1].rstrip("/")
    slide_path = os.path.join(manual_dir, "slide_api.html")
    dialogue_path = os.path.join(manual_dir, "dialogue.json")
    videos_dir = os.path.join(manual_dir, "videos")
    output_path = os.path.join(manual_dir, "manual_with_videos.html")

    if not os.path.exists(slide_path):
        print(f"Error: slide_api.html が見つかりません: {slide_path}")
        sys.exit(1)
    if not os.path.exists(dialogue_path):
        print(f"Error: dialogue.json が見つかりません: {dialogue_path}")
        sys.exit(1)

    with open(slide_path, encoding="utf-8") as f:
        slide_html = f.read()
    with open(dialogue_path, encoding="utf-8") as f:
        dialogue = json.load(f)

    # 存在する section_NN.mp4 のみエントリ化
    video_entries = []
    for idx, section in enumerate(dialogue["sections"]):
        video_file = f"section_{idx:02d}.mp4"
        if os.path.exists(os.path.join(videos_dir, video_file)):
            video_entries.append({
                "index": idx,
                "title": section["section_title"],
                "duration_sec": section["estimated_duration_sec"],
                "file": video_file,
            })

    if not video_entries:
        print(f"Warning: videos/section_*.mp4 が1つも見つかりません")

    video_section = build_video_section(video_entries)

    # <body> の直後に挿入
    m = re.search(r"(<body[^>]*>)", slide_html)
    if m:
        new_html = slide_html[:m.end()] + "\n" + video_section + slide_html[m.end():]
    else:
        # <body>タグが無い異例ケース: 先頭に付加
        new_html = video_section + slide_html

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(new_html)

    print(f"=== 完了 ===")
    print(f"  Output: {output_path}")
    print(f"  Videos embedded: {len(video_entries)}")
    for v in video_entries:
        print(f"    [{v['index']}] {v['title']} ({v['duration_sec']}s)")


if __name__ == "__main__":
    main()
