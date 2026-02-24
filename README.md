# Speak

Real-time speech-to-text for the acme editor using whisper.cpp streaming.

## Features

- Real-time transcription as you speak
- Toggle recording with a single click
- Streams text directly to acme window
- Uses whisper.cpp tiny model for fast inference
- Latency varies by CPU speed (typically a few seconds)

## Requirements

- whisper.cpp with stream example built (https://github.com/ggml-org/whisper.cpp)
- Go 1.21+
- ALSA utils (for audio capture)
- 9p tools for acme integration

## Installation

1. Install and build whisper.cpp:
```bash
git clone https://github.com/ggml-org/whisper.cpp
cd whisper.cpp
make stream
# Install whisper-stream to PATH
cp build/bin/whisper-stream ~/bin/
```

2. Install Speak:
```bash
git clone https://github.com/lneely/acme-speak
cd acme-speak
mk install  # Downloads tiny model and builds binary to ~/bin/Speak
```

## Usage

From acme, add `Speak` to your window tag and click it:
- First click: Start real-time transcription
- Second click: Stop transcription

Text appears in the acme window as you speak.

## Configuration

Edit `main.go` to change:
- Model: Change `modelPath` constant (tiny, base, small, etc.)
- Latency/accuracy: Adjust `--step` and `--length` parameters
  - Current: 500ms step, 5000ms context (good balance)
  - Lower latency: 300ms step, 3000ms context (less accurate)
  - Higher accuracy: 1000ms step, 8000ms context (more lag)

Rebuild after changes: `go build -o ~/bin/Speak .`
