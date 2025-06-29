#!/usr/bin/env -S uv run --script
# cronic:
#   name: Example Python Job
#   desc: Say hello every 5 seconds
#   cron: */5 * * * * *

import sys
import time

print(f"Hello, from example2.py using {sys.executable}")
time.sleep(2)
print(f"Hello, again, from example2.py using {sys.executable}")

