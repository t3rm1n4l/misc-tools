package main

import "net"
import "runtime"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	res := []byte("123456")
	ln, _ := net.Listen("tcp", ":7777")
	for {
		conn, _ := ln.Accept()
		b := make([]byte, 25)
		go func() {
			for {
				_, err := conn.Read(b)
				if err != nil {
					conn.Close()
					return
				}
				conn.Write(res)
			}
		}()
	}

}
