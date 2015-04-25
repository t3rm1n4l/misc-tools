package main

import "strconv"
import "net"
import "os"
import "time"
import "fmt"
import "sync"
import "io"
import "net/http"
import "runtime"
import _ "net/http/pprof"
import "encoding/json"

var (
	ReqSize  = 3
	RespSize = 3
	pool     sync.Pool
)

func stats(w http.ResponseWriter, r *http.Request) {
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	w.WriteHeader(200)
	b, _ := json.Marshal(s)
	w.Write(b)
}

func main() {

	http.HandleFunc("/stats/mem", stats)
	go func() {
		http.ListenAndServe("localhost:9102", nil)
	}()

	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 3)
		},
	}

	if os.Args[1] == "server" {
		s := Server{
			addr: ":7777",
			ch:   make(chan Message),
		}

		go s.Run()

		respBuf := make([]byte, RespSize)
		respBuf = []byte("got")

		for {
			msg := <-s.ch

			go func() {
				//		fmt.Println("server: got reqid", msg.id)
				readBuf := pool.Get().([]byte)
				io.ReadFull(msg.conn, (readBuf))
				pool.Put(readBuf)
				//fmt.Println("server: got req", string(readBuf))
				s.ReleaseReadConn(msg.conn)
				conn := s.GetWriteConn(msg.id)
				//fmt.Println("server: got wr conn for", msg.id)
				conn.Write(respBuf)
				s.ReleaseWriteConn(conn)
				//		fmt.Println("server responded", string(respBuf), "to ", msg.id)
			}()
		}
	} else {

		c := Client{
			n:       1,
			addr:    ":7777",
			connMap: make(map[uint32]chan net.Conn),
		}

		var wg sync.WaitGroup

		reqBuf := make([]byte, ReqSize)
		reqBuf = []byte("doi")

		t0 := time.Now()
		n, _ := strconv.Atoi(os.Args[2])
		for i := 0; i < n; i++ {
			//fmt.Println("client wrote", i)
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				s := c.NewStream()
				conn := c.GetWriteConn(s.id)
				conn.Write(reqBuf)
				//		fmt.Println("client", i, "sent reqid", s.id)
				c.ReleaseWriteConn(conn)

				conn = <-s.connch

				respBuf := pool.Get().([]byte)
				io.ReadFull(conn, respBuf)
				pool.Put(respBuf)

				c.ReleaseReadConn(conn)
				//		fmt.Println("client", i, "got resp", s.id)
			}(&wg)
		}
		wg.Wait()

		fmt.Println(float64(n)/time.Now().Sub(t0).Seconds(), "reqs/sec")
	}
}
