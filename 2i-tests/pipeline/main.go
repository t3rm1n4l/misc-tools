package main

import (
	"encoding/json"
	"fmt"
	c "github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/indexer"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"time"
)

var n = 0

type KV struct {
	k []byte
	v []byte
}

func loader(ch chan *KV) {
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("[\"key-%010d\"]", i)
		v := fmt.Sprintf("doc-%d", i)
		kv := &KV{k: []byte(k), v: []byte(v)}
		ch <- kv
	}
	close(ch)
}

func stats(w http.ResponseWriter, r *http.Request) {
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	w.WriteHeader(200)
	b, _ := json.Marshal(s)
	w.Write(b)
}

func main() {

	http.HandleFunc("/stats/mem", stats)
	go func() {
		http.ListenAndServe("localhost:9102", nil)
	}()

	n, _ = strconv.Atoi(os.Args[1])
	ch := make(chan *KV, 1000)
	go loader(ch)

	conf := c.SystemConfig.SectionConfig("indexer.", true)
	//slice, err := indexer.NewMemDBSlice(indexer.SliceId(0), c.IndexDefnId(0), c.IndexInstId(0), false, nil)
	slice, err := indexer.NewForestDBSlice("/tmp/test.fdb", indexer.SliceId(0), c.IndexDefnId(0), c.IndexInstId(0), false, conf)
	c.CrashOnError(err)

	st := time.Now()
	for kv := range ch {
		err := slice.Insert(kv.k, kv.v)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("loading %v items took %v\n", n, time.Now().Sub(st))

	si, _ := slice.NewSnapshot(nil, false)
	snap, _ := slice.OpenSnapshot(si)
	//var count int = 0

	for i := 0; i < 10; i++ {
		pool := indexer.BlockPool.Get()
		indexer.BlockPool.Put(pool)
	}

	ln, _ := net.Listen("tcp", ":8081")

	for {
		conn, _ := ln.Accept()

		src := indexer.IndexEntrySrc{
			Snapshot: snap,
			Outch:    make(chan *[]byte, 3),
		}

		dec := indexer.IndexEntryDecoder{
			Inch:  src.Outch,
			Outch: make(chan *[]byte, 3),
		}

		send := indexer.IndexEntrySender{
			Inch: dec.Outch,
			W:    conn,
		}

		go src.Routine()
		go dec.Routine()

		send.Routine()

		conn.Close()
	}

	/*
		var r indexer.ByteBlockReader

		st = time.Now()
		for itm := range dec.Outch {
			r.Init(itm)
			for {
				_, err := r.Get()
				if err != nil {
					break
				}

				count++
			}

			indexer.BlockPool.Put(itm)
		}

		fmt.Printf("receiving %d items took %v\n", count, time.Now().Sub(st))
	*/

}
