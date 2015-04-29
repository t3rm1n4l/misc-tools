package main

import (
	"fmt"
	"github.com/couchbase/cbauth"
	"github.com/couchbase/indexing/secondary/common"
	"os"
	"time"
)

func client(id int, server string, cancel bool) {
	sn, err := common.NewServicesChangeNotifier(common.ClusterUrl(server), "default")
	fmt.Printf("client (%v)\n", id)
	if err != nil {
		fmt.Println("error %v", err)
		os.Exit(1)
	}

	//go func() {
	//if cancel {
	//time.Sleep(time.Second * 5)
	//sn.Close()
	//}
	//}()

	for {
		x, err := sn.Get()

		if err != nil {
			if err != nil {
				fmt.Printf("Got error for client %v, %v", id, err)
				os.Exit(1)
			}
		}
		fmt.Printf("Got config for client %v\nconfig:%v\n", id, x.Msg)

	}
}

func main() {
	_, err := cbauth.InternalRetryDefaultInit("localhost:9000", "Administrator", "asdasd")
	if err != nil {
		os.Exit(1)
	}

	for i := 0; i < 1; i++ {
		go client(i, "localhost:9000", false)
	}
	//for i := 20; i < 40; i++ {
	//go client(i, "localhost:9001", false)
	//}
	//for i := 40; i < 60; i++ {
	//go client(i, "localhost:9002", true)
	//}
	time.Sleep(time.Hour)
}
