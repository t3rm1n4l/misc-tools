package main

import (
	"fmt"
	"time"
	)

func main() {
	var arr []string
	N := 10000000
	s := time.Now()
	for i:=0; i < N; i++ {
		arr = append(arr, fmt.Sprintf("document_%07d", i))
	}
	
	fmt.Println("took ", time.Now().Sub(s))
}

	
