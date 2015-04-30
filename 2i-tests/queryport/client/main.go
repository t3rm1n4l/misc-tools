package main

import (
	"fmt"
	c "github.com/couchbase/indexing/secondary/common"
	qclient "github.com/couchbase/indexing/secondary/queryport/client"
	"os"
	"strconv"
	"sync"
	"time"
)

func callb(resp qclient.ResponseReader) bool {
	return true
}

func main() {
	addr := ":7777"
	queryportCfg := c.SystemConfig.SectionConfig("queryport.client.", true)

	var wg sync.WaitGroup
	n, _ := strconv.Atoi(os.Args[1])
	thr, _ := strconv.Atoi(os.Args[2])
	x := queryportCfg["poolSize"]
	x.Value = thr
	queryportCfg["poolSize"] = x
	gsi := qclient.NewGsiScanClient(addr, queryportCfg)

	perthr := n / thr

	t0 := time.Now()
	for i := 0; i < thr; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < perthr; i++ {
				gsi.ScanAll(1, 1, c.AnyConsistency, nil, callb)
			}
		}()
	}

	wg.Wait()

	fmt.Println(float64(n)/time.Now().Sub(t0).Seconds(), "rows/sec")
}
