package main

import (
	"fmt"
	c "github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/indexer"
	"os"
	"strconv"
	"time"
)

var n = 0

type KV struct {
	k indexer.Key
	v indexer.Value
}

func loader(ch chan *KV) {
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("[\"key-%0250d\"]", i)
		v := fmt.Sprintf("doc-%d", i)
		key, _ := indexer.NewKey([]byte(k))
		value, _ := indexer.NewValue([]byte(v))
		kv := &KV{k: key, v: value}
		ch <- kv
	}
	close(ch)
}

func main() {

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
	st = time.Now()
	quitch := make(indexer.StopChannel)
	kch, _ := snap.KeySet(quitch)

	for _ = range kch {
		count++
	}
	fmt.Printf("iterating (ch) %v items took %v\n", n, time.Now().Sub(st))

	st = time.Now()
	count = 0
	var ckr indexer.CallbackKeyReader
	ckr = snap.(indexer.CallbackKeyReader)

	cb := func(k indexer.Key) bool {
		count++
		return true
	}
	ckr.KeySetCb(cb)
	fmt.Printf("iterating (cb) %v items took %v\n", n, time.Now().Sub(st))

}
