package main

import "math/rand"
import "fmt"
import "time"
import "runtime"

func writeSomething(b []byte) {
	for i, _ := range b {
		b[i] = 45
	}
}

func main() {
	refs := [][]byte{}
	max := 5000000000
	sizes := []int{1000000, 10000000, 5000000, 6000000, 500000, 20000}
	alloced := 0
	for {
		i := rand.Int() % len(sizes)
		block := make([]byte, sizes[i], sizes[i])
		writeSomething(block)
		alloced += len(block)
		fmt.Println("alloc", sizes[i], alloced)
		refs = append(refs, block)
		if alloced > max {
			break
		}
	}

	fmt.Println("done")
	time.Sleep(time.Minute)

	fmt.Println("wait for release to OS")
	refs = nil
	runtime.GC()
	time.Sleep(time.Hour)
}
