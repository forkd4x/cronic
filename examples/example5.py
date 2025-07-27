#!/usr/bin/env -S mise exec python@2 -- python
# cronic:
#   name: Example Python2 Job
#   desc: Say hello every 22 seconds
#   cron: */22 * * * * *

import sys
import time

print "Hello, from %s using %s" % (__file__[2:], sys.executable)
time.sleep(11)
print "Bye, from %s using %s" % (__file__[2:], sys.executable)
