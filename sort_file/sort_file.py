#!/usr/bin/env python

import struct
import os

def read_record(f):
    buf = f.read(4)
    if buf == '':
        return None

    rlen = struct.unpack('<I', buf)[0]

    buf = f.read(2)
    klen = struct.unpack('>h', buf)[0]

    vlen = rlen - klen - 2

    key = f.read(klen)
    val = f.read(vlen)

    return key, val

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

print read_records("/Users/sarath/development/couchbase/test_data/index_builder_backup/sort_files/45ec5bb324149fb577a146f0c40028a1.sort", 10)








