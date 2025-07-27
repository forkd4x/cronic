#!/usr/bin/env -S uv run --script
# cronic:
#   name: Example Python Job
#   desc: Say hello every 6 seconds
#   cron: */6 * * * * *

import sys
import time
from pathlib import Path

print(f"Hello, from {Path(__file__).name} using {sys.executable}")
time.sleep(3)
print(f"Bye, from {Path(__file__).name} using {sys.executable}")

