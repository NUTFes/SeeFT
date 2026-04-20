#!/usr/bin/env python3
"""
dialogue.json → ターン別音声WAV + manifest.json

使い方:
  # VOICEVOX バックエンド（推奨、要 localhost:50021 で engine 起動）
  python3 scripts/synthesize_audio.py docs/manuals/01_44th_配線マニュアル --backend voicevox

  # macOS say バックエンド（オフライン、Kyoko 固定）
  python3 scripts/synthesize_audio.py docs/manuals/01_44th_配線マニュアル --backend say

  # 特定セクションのみ
  python3 scripts/synthesize_audio.py docs/manuals/01_44th_配線マニュアル --section 4

VOICEVOX 話者候補:
  Expert  候補: 玄野武宏ノーマル(11), 青山龍星ノーマル(13)
  Novice  候補: 春日部つむぎノーマル(8), 四国めたんノーマル(2), ずんだもんノーマル(3)
"""

import argparse
import json
import os
import re
import subprocess
import sys
import urllib.parse
import urllib.request

VOICEVOX_URL = "http://localhost:50021"

DEFAULT_SAY_VOICE = {"Expert": "Kyoko", "Novice": "Kyoko"}
DEFAULT_SAY_RATE = {"Expert": 175, "Novice": 205}
DEFAULT_VOICEVOX_ID = {"Expert": 11, "Novice": 8}  # 玄野武宏, 春日部つむぎ

PAUSE_SEC = {
    "short pause": 0.25,
    "medium pause": 0.5,
    "long pause": 1.0,
}


def strip_markup(text: str) -> str:
    return re.sub(r"\[[^\]]+\]", "", text).strip()


def synth_say(text: str, voice: str, rate: int, out_wav: str):
    subprocess.run(
        ["say", "-v", voice, "-r", str(rate),
         "--data-format=LEI16@24000",
         "-o", out_wav, text],
        check=True,
    )


def synth_voicevox(text: str, speaker_id: int, out_wav: str):
    """VOICEVOX engine を叩いて wav を生成する (24kHz mono)"""
    # 1. audio_query でプロソディ情報を取得
    q_url = f"{VOICEVOX_URL}/audio_query?speaker={speaker_id}&text={urllib.parse.quote(text)}"
    q_req = urllib.request.Request(q_url, method="POST")
    with urllib.request.urlopen(q_req, timeout=60) as r:
        query = json.loads(r.read())

    # 2. synthesis で wav を生成
    s_url = f"{VOICEVOX_URL}/synthesis?speaker={speaker_id}"
    s_req = urllib.request.Request(
        s_url,
        data=json.dumps(query).encode("utf-8"),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(s_req, timeout=120) as r:
        wav_bytes = r.read()

    with open(out_wav, "wb") as f:
        f.write(wav_bytes)


def get_wav_duration(path: str) -> float:
    r = subprocess.run(
        ["ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", path],
        capture_output=True, text=True, check=True,
    )
    return float(json.loads(r.stdout)["format"]["duration"])


def add_trailing_silence(wav_path: str, silence_sec: float):
    if silence_sec <= 0:
        return
    tmp = wav_path + ".tmp.wav"
    subprocess.run(
        ["ffmpeg", "-y", "-loglevel", "error",
         "-i", wav_path,
         "-af", f"apad=pad_dur={silence_sec}",
         "-c:a", "pcm_s16le", "-ar", "24000", "-ac", "1",
         tmp],
        check=True,
    )
    os.replace(tmp, wav_path)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("manual_dir")
    ap.add_argument("--backend", choices=["say", "voicevox"], default="voicevox")
    ap.add_argument("--section", type=int, help="対象セクションindex (0-based)")
    # say バックエンド用
    ap.add_argument("--expert-voice", default=DEFAULT_SAY_VOICE["Expert"])
    ap.add_argument("--novice-voice", default=DEFAULT_SAY_VOICE["Novice"])
    ap.add_argument("--expert-rate", type=int, default=DEFAULT_SAY_RATE["Expert"])
    ap.add_argument("--novice-rate", type=int, default=DEFAULT_SAY_RATE["Novice"])
    # voicevox バックエンド用
    ap.add_argument("--expert-speaker-id", type=int, default=DEFAULT_VOICEVOX_ID["Expert"])
    ap.add_argument("--novice-speaker-id", type=int, default=DEFAULT_VOICEVOX_ID["Novice"])
    args = ap.parse_args()

    manual_dir = args.manual_dir.rstrip("/")
    dialogue_path = os.path.join(manual_dir, "dialogue.json")
    audio_dir = os.path.join(manual_dir, "audio")
    os.makedirs(audio_dir, exist_ok=True)

    with open(dialogue_path, encoding="utf-8") as f:
        dialogue = json.load(f)

    all_sections = dialogue["sections"]
    target_indices = [args.section] if args.section is not None else list(range(len(all_sections)))

    print(f"=== 音声合成 (backend={args.backend}) ===")
    if args.backend == "voicevox":
        print(f"  Expert speaker_id={args.expert_speaker_id}, Novice speaker_id={args.novice_speaker_id}")
        speaker_of = {"Expert": args.expert_speaker_id, "Novice": args.novice_speaker_id}
    else:
        print(f"  Expert={args.expert_voice}@{args.expert_rate}, Novice={args.novice_voice}@{args.novice_rate}")
        voice_of = {"Expert": args.expert_voice, "Novice": args.novice_voice}
        rate_of = {"Expert": args.expert_rate, "Novice": args.novice_rate}

    manifest = {"manual_name": dialogue["manual_name"], "backend": args.backend, "entries": []}

    for i in target_indices:
        sec = all_sections[i]
        print(f"[section {i}] {sec['section_title']} ({len(sec['turns'])} turns)")
        for j, turn in enumerate(sec["turns"]):
            out_wav = os.path.join(audio_dir, f"sec{i:02d}_turn{j:02d}.wav")
            speaker = turn["speaker"]
            text = strip_markup(turn["text"])

            if args.backend == "voicevox":
                synth_voicevox(text, speaker_of.get(speaker, DEFAULT_VOICEVOX_ID["Novice"]), out_wav)
            else:
                synth_say(text, voice_of.get(speaker, "Kyoko"), rate_of.get(speaker, 185), out_wav)

            trailing_silence = max(
                (PAUSE_SEC.get(tag, 0) for tag in turn.get("tags", [])),
                default=0,
            )
            if trailing_silence:
                add_trailing_silence(out_wav, trailing_silence)

            duration = get_wav_duration(out_wav)
            manifest["entries"].append({
                "section_idx": i,
                "turn_idx": j,
                "audio_path": os.path.relpath(out_wav, manual_dir),
                "duration_sec": round(duration, 3),
                "speaker": speaker,
                "image": turn.get("image"),
                "text": turn["text"],
                "tags": turn.get("tags", []),
            })
            print(f"  sec{i:02d}t{j:02d} {speaker:<7} {duration:5.2f}s img={turn.get('image')}")

    manifest_path = os.path.join(audio_dir, "manifest.json")
    with open(manifest_path, "w", encoding="utf-8") as f:
        json.dump(manifest, f, ensure_ascii=False, indent=2)

    total_dur = sum(e["duration_sec"] for e in manifest["entries"])
    print(f"\n=== 完了 ===")
    print(f"  entries: {len(manifest['entries'])}")
    print(f"  total duration: {total_dur:.1f}s ({int(total_dur)//60}min{int(total_dur)%60:02d}s)")
    print(f"  manifest: {manifest_path}")


if __name__ == "__main__":
    main()
