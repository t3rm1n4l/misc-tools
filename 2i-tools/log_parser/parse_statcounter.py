#!/usr/bin/env python

import sys
import re
import json
from dateutil import parser

prevStats = {}
last_time = None

r = re.compile("^([^ ]*) .*PeriodicStats = ([^}]*})")
if __name__ == '__main__':
    f = open(sys.argv[1])
    counter = sys.argv[2]
    needTs = sys.argv[3]

    for line in f.readlines():
        if "PeriodicStats" in line:
            parsed = r.findall(line)
            ts = parsed[0][0]
            time = parser.parse(ts)
            stats = parsed[0][1]
            stats_json = json.loads(stats)
            if prevStats.has_key(counter):
                last_cntr = prevStats[counter]
                if not stats_json.has_key(counter):
                    continue
                curr_cntr = stats_json[counter]
                d = time-last_time
                secs = d.total_seconds()
                diff = curr_cntr - last_cntr
                if needTs == "yes":
                    print ts, counter, diff, "=", str(diff/secs)+"/sec"
                else:
                    print counter, diff, "=", str(diff/secs)+"/sec"
            prevStats = stats_json
            last_time = time
