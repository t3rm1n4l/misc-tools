#!/usr/bin/env python

import struct
import os
import sys

def read_record(f):
    buf = f.read(4)
    if buf == '':
        return None

    rlen = struct.unpack('<I', buf)[0]

    optype = f.read(1)
    if struct.unpack("b", optype)[0] == 1:
        optype = "remove"
    else:
        optype = "insert"

    buf = f.read(2)
    klen = struct.unpack('>h', buf)[0]

    vlen = rlen - klen - 2 - 1

    key = f.read(klen)
    val = f.read(vlen)

    return optype, key, val

def write_record(f, key, val):
    klen = len(key)
    rlen = 2 + klen + len(val)

    f.write(struct.pack("<I", rlen))
    f.write(struct.pack(">h", klen))
    f.write(key)
    f.write(val)

def read_records(fn, count = 5):
    recs = []
    f = open(fn, 'rb')
    for i in xrange(count):
        rec = read_record(f)
        if not rec:
            break
        recs.append(rec)

    return recs

"""
f = open('/tmp/y', 'wb')
write_record(f, "key1", "val1")
write_record(f, "key2", "val2")
f.close()
"""

if __name__ == '__main__':
    f = sys.argv[1]
    for r in read_records(f, 100000):
        print r








