# /// script
# requires-python = ">=3.12"
# dependencies = ["sonapy"]
#
# [tool.uv.sources]
# sonapy = { path = "../" }
# ///
"""
Streaming transcription with stable timestamps.

Setup:
  wget https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin
  wget https://huggingface.co/ggml-org/whisper-vad/resolve/main/ggml-silero-v6.2.0.bin
  wget https://github.com/thewh1teagle/sona/releases/download/audio-files-v1/synth_5min.wav

Run:
  uv run examples/stable_timestamps.py ggml-tiny.bin ggml-silero-v6.2.0.bin synth_5min.wav
"""

import json
import sys
from sonapy import Sona


def main():
    if len(sys.argv) < 4:
        print(f"Usage: {sys.argv[0]} <model.bin> <vad-model.bin> <audio.wav>")
        sys.exit(1)

    model_path = sys.argv[1]
    vad_model_path = sys.argv[2]
    audio_path = sys.argv[3]

    sona = Sona()
    try:
        sona.load_model(model_path)
        last_start = 0.0

        print("Streaming transcription with stable timestamps:")
        data = {
            "response_format": "json",
            "stream": "true",
            "stable_timestamps": "true",
            "vad_model": vad_model_path,
        }
        with open(audio_path, "rb") as f:
            with sona._client._http.stream(
                "POST",
                "/v1/audio/transcriptions",
                files={"file": (audio_path, f, "application/octet-stream")},
                data=data,
            ) as resp:
                for line in resp.iter_lines():
                    if not line:
                        continue
                    event = json.loads(line)
                    match event["type"]:
                        case "progress":
                            print(f"  progress: {event['progress']}%")
                        case "segment":
                            start = max(float(event["start"]), last_start)
                            last_start = start
                            end = max(float(event.get("end", start)), start)
                            print(f"  [{start:.1f}s-{end:.1f}s] {event['text']}")
                        case "result":
                            pass
    finally:
        sona.stop()


if __name__ == "__main__":
    main()
