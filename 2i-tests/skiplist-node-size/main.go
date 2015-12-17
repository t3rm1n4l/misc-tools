package main

import "fmt"
import "unsafe"
import "sync/atomic"

// ideal
// struct {
// a,b uint32
// ptr *byte
// klen uint16
// next *pointers
// nextLen uint8
//

type ideal struct {
	a, b uint32         // 8
	data unsafe.Pointer // 8

	nextPtr unsafe.Pointer // 8
	n1, n2  unsafe.Pointer // 16

	nextLen uint8  // 1
	dataLen uint16 // 2

}

type Itm struct {
	a, b uint32
	bs   unsafe.Pointer
	l    int
	//bs []byte
}

type Node struct {
	arr []unsafe.Pointer
	//arr unsafe.Pointer
	v    atomic.Value
	next *Node
	//itm  Itm
	//l    uint16
}

type proxy struct {
	n   unsafe.Pointer
	del bool
}

func main() {
	var itm Itm
	var n Node
	var p unsafe.Pointer
	var px proxy

	nodeSz := unsafe.Sizeof(n)
	itmSz := unsafe.Sizeof(itm)
	nextPtrSize := 2 * unsafe.Sizeof(p)
	proxyPtrSize := 2 * unsafe.Sizeof(px)
	fmt.Println("node =", nodeSz, "itm =", itmSz, "next =", nextPtrSize, "proxy =", proxyPtrSize, "=", nodeSz+itmSz+nextPtrSize+proxyPtrSize)

	var nn ideal
	fmt.Println(unsafe.Sizeof(nn))
}
