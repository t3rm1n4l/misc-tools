package main

import "net"
import "encoding/binary"
import "sync"

import "fmt"
import "io"
import "bufio"
import "time"

const bufSize = 4 * 1024

type Server struct {
	sync.Mutex
	addr string

	connPool chan io.Writer
	ch       chan Message
	nr       int
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
		rdr := bufio.NewReaderSize(conn, bufSize)
		wr := bufio.NewWriterSize(conn, bufSize)
		go func() {

			for {
				c.Lock()
				if c.nr == 0 {
					wr.Flush()
				}
				c.Unlock()
				time.Sleep(time.Nanosecond * 1000000)
			}
		}()
		fmt.Println("new serv conn", n)
		c.connPool <- wr
		go c.monitorConn(rdr)
	}
}

func (c *Server) monitorConn(conn io.Reader) {
	var id uint32
	err := binary.Read(conn, binary.LittleEndian, &id)
	if err != nil {
		fmt.Println("conn %s closed", conn)
		c.Lock()
		c.connPool = make(chan io.Writer, 100)
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
	c.Lock()
	c.nr++
	c.Unlock()

	conn := <-c.connPool
	err := binary.Write(conn, binary.LittleEndian, id)
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *Server) ReleaseWriteConn(conn io.Writer) {
	c.Lock()
	defer c.Unlock()
	c.connPool <- conn
	c.nr--
}

func (c *Server) ReleaseReadConn(conn io.Reader) {
	go c.monitorConn(conn)
}

func (c *Server) Flush() {
	close(c.connPool)
	for x := range c.connPool {
		conn := x.(*bufio.Writer)
		conn.Flush()
	}
}
