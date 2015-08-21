package main

import "crypto/rand"
import "github.com/couchbase/go-couchbase"
import "os"
import "strconv"
import "fmt"
import "sync"

import "time"

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func main() {

	c, _ := couchbase.Connect("http://127.0.0.1:9000")
	p, _ := c.GetPool("default")
	b, _ := p.GetBucket("default")

	n, _ := strconv.Atoi(os.Args[1])
	sz, _ := strconv.Atoi(os.Args[2])
	nitr, _ := strconv.Atoi(os.Args[3])
	nthr, _ := strconv.Atoi(os.Args[4])
	docOffset, _ := strconv.Atoi(os.Args[5])

	for x := 0; x < nitr; x++ {
		var wg sync.WaitGroup
		for y := 0; y < nthr; y++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				for i := 0; i < n/nthr; i++ {
					time.Sleep(20 * time.Microsecond)
					docid := fmt.Sprintf("docid-%d", i+offset+docOffset)
					err := b.Set(docid, 0, map[string]interface{}{"name": randString(sz)})
					if err != nil {
						fmt.Println(err)
					}
					//time.Sleep(time.Microsecond * 2000)
				}
			}(y * n / nthr)
		}
		wg.Wait()
	}

}
