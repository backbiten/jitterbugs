# Contributing to Jitterbugs

Thank you for your interest in contributing! This project welcomes issues, bug fixes, and feature suggestions.

## Getting started

1. Fork the repository and clone your fork.
2. Create a virtual environment and install dependencies:
   ```bash
   python -m venv .venv
   source .venv/bin/activate   # Windows: .venv\Scripts\Activate.ps1
   pip install -U pip
   pip install -r requirements.txt
   ```
3. Make your changes on a feature branch.
4. Verify the demo runs locally:
   ```bash
   python scripts/run_sniffer.py
   ```
5. Open a pull request against `main` with a clear description of what you changed and why.

## Code style

- Follow [PEP 8](https://peps.python.org/pep-0008/).
- Keep functions short and well-documented with docstrings.
- No secrets, credentials, or personal data in the repository.

## Scope

Please keep contributions in scope:

- Real-time audio/video **feature extraction and visualisation** only.
- Consent-based, local processing.
- No surveillance, covert capture, or any system that infers criminality or intent.

See `CODE_OF_CONDUCT.md` for our community standards.
