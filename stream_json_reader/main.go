package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {
	url := "http://localhost:9000/pools/default/nodeServicesStreaming"
	resp, _ := http.Get(url)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		fmt.Println(string(line))
	}
}
