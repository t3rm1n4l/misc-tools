package main

import "net"
import "sync"
import "encoding/binary"
import "runtime"
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
	connPool []io.Writer
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
	c.Lock()
	defer c.Unlock()
	c.connPool = append(c.connPool, conn)
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

	rdr := bufio.NewReaderSize(conn, 1024*4)
	wr := bufio.NewWriterSize(conn, 1024*4)
	c.connPool = append(c.connPool, wr)
	go c.monitorConn(rdr)
}

func (c *Client) GetWriteConn(id uint32) io.Writer {
retry:
	c.Lock()
	if len(c.connPool) == 0 {
		if c.conns < maxConns {
			c.addConn()
			c.conns++
		} else {
			c.Unlock()
			runtime.Gosched()
			goto retry
		}
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

func (c *Client) Close() {
	conn := c.connPool[0]
	conn.(*bufio.Writer).Flush()
}
