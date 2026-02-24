# Speak

Speech-to-text tool for the acme editor using whisper.cpp streaming.

## Requirements

- whisper.cpp with stream example built (https://github.com/ggml-org/whisper.cpp)
- ALSA utils (for audio capture)
- 9p tools for acme integration

## Installation

1. Install whisper.cpp and build the stream example:
```bash
cd whisper.cpp
make stream
# Install to PATH
cp build/bin/whisper-stream ~/bin/
```

2. Download a model:
```bash
cd whisper.cpp/models
./download-ggml-model.sh tiny
```

3. Set WHISPER_MODEL environment variable or edit the Speak script

4. Install Speak:
```bash
mk install
```

## Usage

From acme, add `Speak` to your window tag and click it:
- First click: Start real-time transcription (streams text as you speak)
- Second click: Stop transcription

The script uses the `$winid` environment variable set by acme to write transcribed text to the current window in real-time.
