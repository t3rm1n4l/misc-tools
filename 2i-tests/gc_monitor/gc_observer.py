import requests
import json
import time

last = None

gct = None

def loop():
    global last
    global gct

    while True:
        r = requests.get('http://localhost:9102/stats/mem')
        itm = json.loads(r.content)

        if last:
            print "AllocsPerSec", (itm["TotalAlloc"]-last["TotalAlloc"])/(1024*1024.0)
            print "HeapInuse", itm["HeapInuse"] / (1024*1024.0)
            print "AvgGCpause", (sum(itm["PauseNs"])/(len(itm["PauseNs"])))/(10**9 * 1.000)
            print "NumGCs", itm["NumGC"]

            gcint = itm["LastGC"] - last["LastGC"]
            if gcint:
                gct = gcint/(10**9 * 1.0)
            print "GCInterval", gct
            print "NumHeapObjects", itm["HeapObjects"]
            print "NumHeapObjectsPerSec", (itm["HeapObjects"] - last["HeapObjects"])
            print "NextGC", itm["NextGC"]/(1024*1024.0)

            print

        last = itm
        time.sleep(1)


loop()
