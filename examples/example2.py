#!/usr/bin/env -S uv run --script
# cronic:
#   name: Example Python Job
#   desc: Say hello every 5 seconds
#   cron: */5 0 0 0 0 0
import sys

print(f"Hello, from example2.py using {sys.executable}")
