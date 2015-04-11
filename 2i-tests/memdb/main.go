package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

type itm string

func (i itm) Less(than Item) bool {
	return i < *than.(*itm)
}

func main() {
	l := NewLLRB()

	ch := make(chan *itm, 10000)

	n, _ := strconv.Atoi(os.Args[1])

	go func() {
		for i := 0; i < n; i++ {
			x := itm(fmt.Sprintf("itm-%d", i))
			ch <- &x
		}
		close(ch)
	}()

	prof, _ := os.Create("prof")
	pprof.StartCPUProfile(prof)
	defer pprof.StopCPUProfile()
	t := time.Now()
	//for x := range ch {
	//l.InsertNoReplace(x)
	//}

	f, _ := bolt.Open("/tmp/bolt.db", 777, nil)
	txn := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("default"))
		return err
	}
	f.Batch(txn)

	f.NoSync = true
	txn = func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("default"))
		for x := range ch {
			b.Put([]byte(*x), nil)
		}
		return nil
	}

	f.Update(txn)

	fmt.Printf("loading %v items took %v\n", n, time.Now().Sub(t))

	t = time.Now()
	n = count(l.root)
	fmt.Printf("counting %v items took %v\n", n, time.Now().Sub(t))

}
