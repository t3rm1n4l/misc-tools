package main

import (
	"fmt"
	c "github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/indexer"
	"github.com/couchbase/indexing/secondary/queryport"
	"net"
)

func serverCallback(req interface{}, conn net.Conn, q <-chan interface{}) {
	w := indexer.NewProtoWriter(indexer.ScanAllReq, conn)
	w.Row([]byte("pk"), nil)
	w.Done()
}

func main() {
	addr := ":7777"
	ch := make(chan int)
	queryportCfg := c.SystemConfig.SectionConfig("indexer.queryport.", true)
	_, err := queryport.NewServer(addr, serverCallback, queryportCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-ch
}
