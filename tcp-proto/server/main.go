package main

import "net"
import "runtime"
import "strings"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	res := make([]byte, 200)
	ln, _ := net.Listen("tcp", ":7777")
	for {
		newconn, _ := ln.Accept()
		//conn.(*net.TCPConn).SetWriteBuffer(16 * 1024)
		//conn.(*net.TCPConn).SetReadBuffer(16 * 1024)

		go func(conn net.Conn) {
		loop:
			for { // request

				b := make([]byte, 200)
				_, err := conn.Read(b)
				if err != nil {
					//conn.Close()
					if strings.Contains(err.Error(), "EOF") {
						break loop
					}
				}

				ch1 := make(chan bool)
				ch2 := make(chan bool)
				donech := make(chan bool)

				go func() {
					<-ch1
					go func() {
						select {
						case <-ch1:
						case <-donech:
						}
					}()
					conn.Write(res)
					donech <- true
					ch2 <- true
				}()
				ch1 <- true
				<-ch2
			}
		}(newconn)
	}

}
