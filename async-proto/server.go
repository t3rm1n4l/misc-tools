package main

import "net"
import "encoding/binary"
import "sync"

import "fmt"
import "runtime"

type Server struct {
	sync.Mutex
	addr string

	connPool []net.Conn
	ch       chan Message
}

type Message struct {
	id   uint32
	conn net.Conn
}

func (c *Server) Run() {
	ln, _ := net.Listen("tcp", c.addr)
	n := 0
	for {
		conn, _ := ln.Accept()
		n++
		fmt.Println("new serv conn", n)
		conn.(*net.TCPConn).SetWriteBuffer(1024 * 4)
		conn.(*net.TCPConn).SetReadBuffer(1024 * 4)
		c.Lock()
		c.connPool = append(c.connPool, conn)
		c.Unlock()
		go c.monitorConn(conn)
	}
}

func (c *Server) monitorConn(conn net.Conn) {
	var id uint32
	err := binary.Read(conn, binary.LittleEndian, &id)
	if err != nil {
		fmt.Println("conn %s closed", conn)
		c.Lock()
		c.connPool = make([]net.Conn, 0)
		c.Unlock()
		return
	}
	//	fmt.Println("server got conn for ", id)
	c.ch <- Message{id: id, conn: conn}
}

func (s *Server) Channel() chan Message {
	return s.ch
}

func (c *Server) GetWriteConn(id uint32) net.Conn {

retry:
	c.Lock()
	if len(c.connPool) == 0 {
		c.Unlock()
		runtime.Gosched()
		goto retry
	}

	defer c.Unlock()
	conn := c.connPool[0]
	err := binary.Write(conn, binary.LittleEndian, id)
	if err != nil {
		panic(err)
	}
	c.connPool = c.connPool[1:]
	return conn
}

func (c *Server) ReleaseWriteConn(conn net.Conn) {
	c.Lock()
	defer c.Unlock()
	c.connPool = append(c.connPool, conn)
}

func (c *Server) ReleaseReadConn(conn net.Conn) {
	go c.monitorConn(conn)
}
