"""Entry point: run the Jitterbugs overlay demo on camera index 0."""

import os
import sys

# Path manipulation must precede the package import so the source tree is
# found when running the script directly without installing the package.
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "src"))

from jitterbugs.app import run  # noqa: E402

if __name__ == "__main__":
    run(camera_index=0)
