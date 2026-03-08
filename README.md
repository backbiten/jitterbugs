# jitterbugs

Cross-platform (Windows/macOS/Linux) real-time webcam + microphone feature overlay.

This repo is intended as an open-source reference implementation for **consented** audio/vision feature extraction and visualization:
- Pose (shoulder tilt) via MediaPipe
- Face landmarks (simple gaze proxy + rough eye openness ratio) via MediaPipe FaceMesh
- Audio features (pitch + jitter proxy + RMS energy) via PyAudio + librosa
- On-screen overlay only (no logging)

Press `q` to quit.

## What this is (and isn’t)
This project is a **feature extraction + visualization** demo. It does **not** claim to detect deception, intent, or personal attributes.

## Requirements
- Python 3.10 or 3.11 recommended
- Webcam + microphone

## Install
```bash
python -m venv .venv
```

Activate:

**Windows (PowerShell)**
```bash
\..venv\Scripts\Activate.ps1
```

**macOS/Linux**
```bash
source .venv/bin/activate
```

Install:
```bash
python -m pip install -U pip
pip install -r requirements.txt
```

### PyAudio platform notes
**macOS**
```bash
brew install portaudio
pip install -r requirements.txt
```

**Linux (Debian/Ubuntu)**
```bash
sudo apt-get update
sudo apt-get install -y portaudio19-dev
pip install -r requirements.txt
```

## Run
```bash
python scripts/run_sniffer.py
```