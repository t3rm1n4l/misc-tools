package main

import "net"
import "sync"
import "encoding/binary"
import "runtime"

import "fmt"

const (
	maxConns = 1
)

type Client struct {
	addr string
	sync.Mutex
	n        uint32
	connMap  map[uint32]chan net.Conn
	connPool []net.Conn
	conns    int
}

type Stream struct {
	id     uint32
	connch chan net.Conn
}

func (c *Client) NewStream() Stream {
	c.Lock()
	defer c.Unlock()

	c.n++
	ch := make(chan net.Conn, 1)
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

func (c *Client) ReleaseReadConn(conn net.Conn) {
	go c.monitorConn(conn)
}

func (c *Client) ReleaseWriteConn(conn net.Conn) {
	c.Lock()
	defer c.Unlock()
	c.connPool = append(c.connPool, conn)
}

func (c *Client) monitorConn(conn net.Conn) {
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
	conn.SetWriteBuffer(1024 * 4)
	conn.SetReadBuffer(1024 * 4)
	if err != nil {
		panic(err)
	}

	c.connPool = append(c.connPool, conn)
	go c.monitorConn(conn)
}

func (c *Client) GetWriteConn(id uint32) net.Conn {
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
