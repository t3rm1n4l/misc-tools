package main

import "net"
import "os"
import "fmt"
import "time"
import "strconv"
import "io"
import "sync"
import "runtime"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	req := make([]byte, 50)
	n, _ := strconv.Atoi(os.Args[1])
	total := n
	nconns, _ := strconv.Atoi(os.Args[2])
	perconn := n / nconns

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0:7777")

	worker := func(wg *sync.WaitGroup, conn net.Conn, n int) {
		defer wg.Done()
		resp := make([]byte, 50)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		i := 0
		for i = 1; i <= n; i++ {
			_, err = conn.Write(req)
			if err != nil {
				fmt.Println(err)
				break
			}

			io.ReadFull(conn, resp)
		}
	}

	for i := 0; i < nconns; i++ {
		var x int
		wg.Add(1)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if i == nconns-1 {
			x = n
		} else {
			x = perconn
			n -= perconn
		}

		go worker(&wg, conn, x)
	}
	t0 := time.Now()
	wg.Wait()
	dur := time.Now().Sub(t0)
	fmt.Println(dur)

	fmt.Printf("%d rows/sec\n", int(float64(total)/dur.Seconds()))

}
