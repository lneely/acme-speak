# Speak

Speech-to-text tool for the acme editor using whisper.cpp.

## Requirements

- whisper.cpp (https://github.com/ggml-org/whisper.cpp)
- arecord (ALSA utils)
- 9p tools for acme integration

## Installation

1. Install whisper.cpp and download a model (e.g., base.en)

2. Edit the `Speak` script to set paths:
```bash
WHISPER_DIR="/path/to/whisper.cpp"
MODEL="$WHISPER_DIR/models/ggml-base.en.bin"
```

3. Make script executable:
```bash
chmod +x Speak
```

## Usage

From acme, add `Speak` to your window tag and click it:
- First click: Start recording
- Second click: Stop recording and transcribe to current window

The script uses the `$winid` environment variable set by acme to write transcribed text to the current window.
