package main

import (
	"fmt"
	c "github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/indexer"
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
		k := fmt.Sprintf("[\"key-%0250d\"]", i)
		v := fmt.Sprintf("doc-%d", i)
		kv := &KV{k: []byte(k), v: []byte(v)}
		ch <- kv
	}
	close(ch)
}

func main() {

	var stats runtime.MemStats
	var pause uint64

	n, _ = strconv.Atoi(os.Args[1])
	ch := make(chan *KV, 1000)
	go loader(ch)

	slice, err := indexer.NewMemDBSlice(indexer.SliceId(0), c.IndexDefnId(0), c.IndexInstId(0), false, nil)
	c.CrashOnError(err)

	st := time.Now()
	for kv := range ch {
		slice.Insert(kv.k, kv.v)
	}
	fmt.Printf("loading %v items took %v\n", n, time.Now().Sub(st))

	si, _ := slice.NewSnapshot(nil, false)
	snap, _ := slice.OpenSnapshot(si)
	var count int = 0

	runtime.ReadMemStats(&stats)
	pause = stats.PauseTotalNs
	st = time.Now()
	quitch := make(indexer.StopChannel)
	kch, _ := snap.KeySet(quitch)

	for _ = range kch {
		count++
	}
	runtime.ReadMemStats(&stats)
	fmt.Println(stats)
	fmt.Printf("iterating (ch) %v items took %v\n", n, time.Now().Sub(st))
	fmt.Printf("GC pause = %v ns\n", stats.PauseTotalNs-pause)
	pause = stats.PauseTotalNs

	st = time.Now()
	count = 0
	var ckr indexer.CallbackKeyReader
	ckr = snap.(indexer.CallbackKeyReader)

	cb := func(k indexer.Key) bool {
		count++
		return true
	}
	ckr.KeySetCb(cb)
	runtime.ReadMemStats(&stats)
	fmt.Printf("iterating (cb) %v items took %v\n", n, time.Now().Sub(st))
	fmt.Printf("GC pause = %v ns\n", stats.PauseTotalNs-pause)

}
