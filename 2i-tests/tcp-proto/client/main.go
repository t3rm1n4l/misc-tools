package main

import "net"
import "os"
import "fmt"
import "time"
import "strconv"
import "sync"

func main() {

	var wg sync.WaitGroup
	req := make([]byte, 25)

	n, _ := strconv.Atoi(os.Args[1])
	nconns, _ := strconv.Atoi(os.Args[2])

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0:7777")

	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()
		resp := make([]byte, 6)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		i := 0
		for i = 1; i <= n; i++ {
			_, err := conn.Write(req)
			if err != nil {
				fmt.Println(err)
				break
			}
			_, err = conn.Read(resp)
		}
	}

	t0 := time.Now()
	for i := 0; i < nconns; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	dur := time.Now().Sub(t0)

	fmt.Printf("%d rows/sec\n", int(float64(n*nconns)/dur.Seconds()))

}
