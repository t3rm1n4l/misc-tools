package main

import "strconv"
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
	ReqSize  = 100
	RespSize = 100
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

	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, ReqSize)
		},
	}

	if os.Args[1] == "server" {
		s := Server{
			addr:     "localhost:7777",
			ch:       make(chan Message),
			connPool: make(chan io.Writer, 100),
		}
		go func() {
			http.ListenAndServe("localhost:9102", nil)
		}()

		go s.Run()

		respBuf := make([]byte, RespSize)

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
		go func() {
			http.ListenAndServe("localhost:9103", nil)
		}()

		c := Client{
			n:        0,
			addr:     "localhost:7777",
			connMap:  make(map[uint32]chan io.Reader),
			connPool: make(chan io.Writer, maxConns),
		}

		var wg1 sync.WaitGroup

		reqBuf := make([]byte, ReqSize)

		t0 := time.Now()
		n, _ := strconv.Atoi(os.Args[2])
		threads, _ := strconv.Atoi(os.Args[3])
		perthread := n / threads
		for i := 0; i < threads; i++ {
			//fmt.Println("client wrote", i)
			wg1.Add(1)
			go func() {
				defer wg1.Done()

				for x := 0; x < perthread; x++ {
					s := c.NewStream()
					conn := c.GetWriteConn(s.id)
					conn.Write(reqBuf)
					//				fmt.Println("client", i, "sent reqid", s.id)
					c.ReleaseWriteConn(conn)

					conn2 := <-s.connch

					respBuf := pool.Get().([]byte)
					io.ReadFull(conn2, respBuf)

					c.ReleaseReadConn(conn2)
					//				fmt.Println("client", i, "got resp", s.id, string(respBuf))
					pool.Put(respBuf)
					c.CloseStream(s.id)
				}
			}()
		}
		wg1.Wait()

		fmt.Println(float64(n)/time.Now().Sub(t0).Seconds(), "reqs/sec")
	}
}
