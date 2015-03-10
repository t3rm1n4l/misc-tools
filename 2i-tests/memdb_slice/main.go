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

	quitch := make(indexer.StopChannel)
	si, _ := slice.NewSnapshot(nil, false)
	snap, _ := slice.OpenSnapshot(si)
	kch, _ := snap.KeySet(quitch)

	st = time.Now()
	var count int = 0
	for _ = range kch {
		count++
	}
	fmt.Printf("iterating %v items took %v\n", n, time.Now().Sub(st))

}
