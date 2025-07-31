#!/usr/bin/env -S uv run --script
# cronic:
#   name: Example Python Job
#   desc: Say hello every 20 seconds
#   cron: */20 * * * * *

import sys
import time
from pathlib import Path

print(f"Hello, from {Path(__file__).name} using {sys.executable}")
time.sleep(10)
print(f"Bye, from {Path(__file__).name} using {sys.executable}")

