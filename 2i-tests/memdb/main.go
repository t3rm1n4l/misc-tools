package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/couchbase/goforestdb"
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
	o, _ := strconv.Atoi(os.Args[2])
	typ := os.Args[3]

	go func() {
		for i := o; i < n; i++ {
			x := itm(fmt.Sprintf("itm-%0256d", i))
			ch <- &x
		}
		close(ch)
	}()

	prof, _ := os.Create("prof")
	pprof.StartCPUProfile(prof)
	defer pprof.StopCPUProfile()
	t := time.Now()
	if typ == "memdb" {
		for x := range ch {
			l.InsertNoReplace(x)
		}
	} else if typ == "bolt" {

		f, _ := bolt.Open("/tmp/bolt.db", 0777, nil)
		txn := func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("default"))
			return err
		}
		f.Batch(txn)

		f.NoSync = true
		done := false

		for !done {
			txn = func(tx *bolt.Tx) error {
				sz := 250 * 1024 * 1024
				b := tx.Bucket([]byte("default"))
				for x := range ch {
					sz -= len(*x)
					b.Put([]byte(*x), nil)
					if sz <= 0 {
						return nil
					}
				}
				done = true
				return nil
			}

			f.Update(txn)
		}
	} else if typ == "forestdb" {
		//memQuota := uint64(1024 * 1024 * 25)
		//config := forestdb.DefaultConfig()
		//config.SetDurabilityOpt(forestdb.DRB_ASYNC)
		//config.SetBufferCacheSize(memQuota)
		dbfile, err := forestdb.Open("/tmp/forest.db", nil)
		fmt.Println(err)
		kvstore, _ := dbfile.OpenKVStoreDefault(nil)
		for x := range ch {
			kvstore.SetKV([]byte(*x), nil)
		}
		//dbfile.Commit(forestdb.COMMIT_MANUAL_WAL_FLUSH)
		kvstore.Close()
		dbfile.Close()
	}

	fmt.Printf("loading %v items took %v\n", n, time.Now().Sub(t))

	t = time.Now()
	n = count(l.root)
	fmt.Printf("counting %v items took %v\n", n, time.Now().Sub(t))

}
