package main

import "net"
import "sync"
import "encoding/binary"
import "io"
import "bufio"

import "fmt"

const (
	maxConns = 1
)

type Client struct {
	addr string
	sync.Mutex
	n        uint32
	connMap  map[uint32]chan io.Reader
	connPool chan io.Writer
	conns    int
}

type Stream struct {
	id     uint32
	connch chan io.Reader
}

func (c *Client) NewStream() Stream {
	c.Lock()
	defer c.Unlock()

	c.n++
	ch := make(chan io.Reader, 1)
	c.connMap[c.n] = ch

	return Stream{
		id:     c.n,
		connch: ch,
	}
}

func (c *Client) CloseStream(id uint32) {
	c.Lock()
	defer c.Unlock()

	delete(c.connMap, id)
}

func (c *Client) ReleaseReadConn(conn io.Reader) {
	go c.monitorConn(conn)
}

func (c *Client) ReleaseWriteConn(conn io.Writer) {
	c.connPool <- conn
}

func (c *Client) monitorConn(conn io.Reader) {
	var id uint32
	err := binary.Read(conn, binary.LittleEndian, &id)
	if err != nil {
		panic(err)
	}

	c.Lock()
	defer c.Unlock()
	ch := c.connMap[id]

	ch <- conn
}

func (c *Client) addConn() {
	fmt.Println("new conn")
	tcpAddr, _ := net.ResolveTCPAddr("tcp", c.addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	rdr := bufio.NewReaderSize(conn, bufSize)
	wr := bufio.NewWriterSize(conn, bufSize)
	c.connPool <- wr
	go c.monitorConn(rdr)
}

func (c *Client) GetWriteConn(id uint32) io.Writer {
	c.Lock()
	defer c.Unlock()

	if c.conns < maxConns {
		c.addConn()
		c.conns++
	}

	conn := <-c.connPool
	err := binary.Write(conn, binary.LittleEndian, id)
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *Client) Close() {
	close(c.connPool)
	for conn := range c.connPool {
		conn.(*bufio.Writer).Flush()
	}
}
