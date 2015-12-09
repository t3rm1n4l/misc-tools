package main

import "fmt"
import "sync/atomic"
import "unsafe"

type Itm struct {
	a, b uint32
	bs   []byte
}

type Node struct {
	arr  []unsafe.Pointer
	v    atomic.Value
	next *Node
}

type proxy struct {
	n   *Node
	del bool
}

func main() {
	var itm Itm
	var n Node
	var p unsafe.Pointer
	var px proxy
	fmt.Println(unsafe.Sizeof(n) + unsafe.Sizeof(itm) + 2*(unsafe.Sizeof(p)+unsafe.Sizeof(px)))
}
