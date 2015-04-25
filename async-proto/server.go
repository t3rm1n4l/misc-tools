package main

import "net"
import "encoding/binary"
import "sync"

import "fmt"
import "runtime"
import "io"
import "bufio"
import "time"

type Server struct {
	sync.Mutex
	addr string

	connPool []io.Writer
	ch       chan Message
}

type Message struct {
	id   uint32
	conn io.Reader
}

func (c *Server) Run() {
	ln, _ := net.Listen("tcp", c.addr)
	n := 0
	for {
		conn, _ := ln.Accept()
		n++
		rdr := bufio.NewReaderSize(conn, 1024*4)
		wr := bufio.NewWriterSize(conn, 1024*4)
		go func() {
			for {
				time.Sleep(time.Nanosecond * 1000000)
				wr.Flush()
			}
		}()
		fmt.Println("new serv conn", n)
		c.Lock()
		c.connPool = append(c.connPool, wr)
		c.Unlock()
		go c.monitorConn(rdr)
	}
}

func (c *Server) monitorConn(conn io.Reader) {
	var id uint32
	err := binary.Read(conn, binary.LittleEndian, &id)
	if err != nil {
		fmt.Println("conn %s closed", conn)
		c.Lock()
		c.connPool = make([]io.Writer, 0)
		c.Unlock()
		return
	}
	//	fmt.Println("server got conn for ", id)
	c.ch <- Message{id: id, conn: conn}
}

func (s *Server) Channel() chan Message {
	return s.ch
}

func (c *Server) GetWriteConn(id uint32) io.Writer {

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

func (c *Server) ReleaseWriteConn(conn io.Writer) {
	c.Lock()
	defer c.Unlock()
	c.connPool = append(c.connPool, conn)
}

func (c *Server) ReleaseReadConn(conn io.Reader) {
	go c.monitorConn(conn)
}

func (c *Server) Flush() {
	conn := c.connPool[0].(*bufio.Writer)
	conn.Flush()
}
