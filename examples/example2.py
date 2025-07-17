#!/usr/bin/env -S uv run --script
# cronic:
#   name: Example Python Job
#   desc: Say hello every 6 seconds
#   cron: */6 * * * * *

import sys
import time

print(f"Hello, from example2.py using {sys.executable}")
time.sleep(3)
print(f"Bye, from example2.py using {sys.executable}")

