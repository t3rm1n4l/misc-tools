package main

import "net"
import "os"
import "fmt"
import "time"
import "strconv"
import "io"
import "sync"
import "runtime"
import "sync/atomic"

var bytes uint64

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	n, _ := strconv.Atoi(os.Args[1])
	total := n
	nconns, _ := strconv.Atoi(os.Args[2])
	perconn := n / nconns

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0:7777")
	var conns []net.Conn
	for i := 0; i < nconns; i++ {
		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		conns = append(conns, conn)
	}

	t0 := time.Now()

	worker := func(wg *sync.WaitGroup, conn net.Conn, n int) {
		defer wg.Done()
		resp := make([]byte, 200)
		/*conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		*/
		var err error
		i := 0
		for i = 1; i <= n; i++ {
			req := make([]byte, 200)
			_, err = conn.Write(req)
			if err != nil {
				fmt.Println(err)
				break
			}

			y, err := io.ReadFull(conn, resp)
			if err != nil {
				panic(err)
			}
			atomic.AddUint64(&bytes, uint64(y))
		}
	}

	for i := 0; i < nconns; i++ {
		var x int
		wg.Add(1)

		if i == nconns-1 {
			x = n
		} else {
			x = perconn
			n -= perconn
		}

		go worker(&wg, conns[i], x)
	}
	wg.Wait()
	dur := time.Now().Sub(t0)
	fmt.Println(dur)

	fmt.Printf("%d rows/sec\n", int(float64(total)/dur.Seconds()))
	fmt.Println("bytes=", bytes)

}
