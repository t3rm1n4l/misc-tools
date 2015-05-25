import requests
import sys
import time
import json
import plotly.plotly as py
from plotly.graph_objs import *


bucket = sys.argv[1]
index = sys.argv[2]
timeout = int(sys.argv[3])

flush_queue_size_lbl = "%s:%s:flush_queue_size" %(bucket, index)
num_flush_queued_lbl = "%s:%s:num_flush_queued" %(bucket, index)
num_docs_indexed_lbl = "%s:%s:num_docs_indexed" %(bucket, index)
num_docs_pending_lbl = "%s:%s:num_docs_pending" %(bucket, index)
num_docs_queued_lbl = "%s:%s:num_docs_queued" %(bucket, index)

num_mutations_queued_lbl = "%s:num_mutations_queued" %(bucket)
mutations_queue_size_lbl = "%s:mutation_queue_size" %(bucket)

vals = {flush_queue_size_lbl:[], num_flush_queued_lbl:[], num_docs_indexed_lbl:[], num_docs_pending_lbl:[],
        num_docs_queued_lbl:[], num_mutations_queued_lbl:[], mutations_queue_size_lbl:[]}

start = time.time()
tvals = []
while True:
    now = time.time()
    elapsed = now-start
    if elapsed > timeout:
        break
    tvals.append(elapsed)
    r = requests.get('http://localhost:9102/stats?async=false')
    stats = json.loads(r.content)
    for k in vals:
        vals[k].append(stats[k])

    time.sleep(1)

flush_queue_size = Scatter(
        name="flush_queue_size",
        x = tvals,
        y=vals[flush_queue_size_lbl])

mutations_queue_size = Scatter(
        name = "mutations_queue_size",
        x = tvals,
        y=vals[mutations_queue_size_lbl])

data1 = Data([flush_queue_size, mutations_queue_size])
py.plot(data1, filename="queue_size")

num_flush_queued = Scatter(
    name="num_flush_queued",
    x=tvals,
    y=vals[num_flush_queued_lbl])

num_mutations_queued = Scatter(
    name="num_mutations_queued",
    x=tvals,
    y=vals[num_mutations_queued_lbl])

data2 = Data([num_flush_queued, num_mutations_queued])

py.plot(data2, filename="items_queued")

