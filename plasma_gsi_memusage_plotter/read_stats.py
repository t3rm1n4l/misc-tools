import subprocess
import requests
import json
import time

def get_rss():
    for line in subprocess.check_output(["ps", "-eo", "comm,rss"]).split("\n"):
        if "indexer" in line:
            return int(line.split(" ")[-1])*1024

    return 0

def get_storage_size():
    try:
        r = requests.get("http://localhost:9102/stats/storage")
        storage_stats = json.loads(r.text)
        return storage_stats[0]["Stats"]["MainStore"]["memory_size"] + storage_stats[0]["Stats"]["BackStore"]["memory_size"]
    except:
        return 0

def num_records():
    try:
        r = requests.get("http://localhost:9102/stats/storage")
        storage_stats = json.loads(r.text)
        mrecs = int(storage_stats[0]["Stats"]["MainStore"]["num_rec_allocs"]) - int(storage_stats[0]["Stats"]["MainStore"]["num_rec_frees"]) + int(storage_stats[0]["Stats"]["MainStore"]["num_rec_swapout"]) - int(storage_stats[0]["Stats"]["MainStore"]["num_rec_swapin"])

        brecs = int(storage_stats[0]["Stats"]["BackStore"]["num_rec_allocs"]) - int(storage_stats[0]["Stats"]["BackStore"]["num_rec_frees"]) + int(storage_stats[0]["Stats"]["BackStore"]["num_rec_swapout"]) - int(storage_stats[0]["Stats"]["BackStore"]["num_rec_swapin"])
        return mrecs,brecs
    except:
        return 0,0

def get_queue_size():
    try:
        r = requests.get("http://localhost:9102/stats")
        return int(json.loads(r.text)["memory_used_queue"])
    except:
        return 0

def get_storage_pending():
    try:
        r = requests.get("http://localhost:9102/stats/storage")
        storage_stats = json.loads(r.text)
        return storage_stats[0]["Stats"]["MainStore"]["reclaim_pending"] + storage_stats[0]["Stats"]["BackStore"]["reclaim_pending"]
    except:
        return 0

def get_jemalloc():
    try:
        r = requests.get("http://localhost:9102/stats/storage/mm")
        for line in r.text.split("\n"):
            if "resident" in line:
                a = line.split(",")
                return int(a[0].split(":")[-1]), int(a[3].split(":")[-1])
    except:
        return 0,0


print "TIME,RSS,STORAGE_USED,JE_ALLOCATED,JE_RESIDENT,RECLAIM_PENDING,QUEUE_SIZE,MAINRECORD,DOCRECORDS"

i = 0
while True:
	allocated, resident = get_jemalloc()
        mrecs,brecs = num_records()
	print "{},{},{},{},{},{},{},{},{}".format(i, get_rss(),get_storage_size(), allocated, resident, get_storage_pending(), get_queue_size(),mrecs,brecs)
        time.sleep(1)
        i+=1
