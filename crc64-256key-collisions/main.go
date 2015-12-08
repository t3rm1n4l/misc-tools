package main

import "hash/crc64"
import "fmt"
import "strconv"
import "os"

func main() {

	var collisions int
	m := make(map[uint64]int)
	tab := crc64.MakeTable(crc64.ISO)

	n, _ := strconv.Atoi(os.Args[1])

	for i := 0; i < n; i++ {
		k := []byte(fmt.Sprintf("doc-%d", i))
		h := crc64.Checksum(k, tab)
		if v, ok := m[h]; !ok {
			m[h] = 1
		} else {
			m[h] = v + 1
		}

	}

	for _, v := range m {
		if v != 1 {
			collisions++
		}
	}

	fmt.Println("collisions", collisions)

}
