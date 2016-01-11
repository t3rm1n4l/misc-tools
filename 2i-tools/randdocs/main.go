package main

import "crypto/rand"
import "github.com/couchbase/go-couchbase"
import "os"
import "strconv"
import "fmt"
import "sync"
import "runtime"

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

	if len(os.Args) != 9 {
		fmt.Println("./randdocs cluster:port bucket count field_sz junk_field_sz iterations threads doc_offset")
		os.Exit(1)
	}

	cluster := os.Args[1]
	c, _ := couchbase.Connect(fmt.Sprintf("http://%s", cluster))
	p, _ := c.GetPool("default")
	b, _ := p.GetBucket(os.Args[2])

	n, _ := strconv.Atoi(os.Args[3])
	sz, _ := strconv.Atoi(os.Args[4])
	junkSz, _ := strconv.Atoi(os.Args[5])
	nitr, _ := strconv.Atoi(os.Args[6])
	nthr, _ := strconv.Atoi(os.Args[7])
	docOffset, _ := strconv.Atoi(os.Args[8])

	runtime.GOMAXPROCS(nthr)

	for x := 0; x < nitr; x++ {
		var wg sync.WaitGroup
		for y := 0; y < nthr; y++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				for i := 0; i < n/nthr; i++ {
					time.Sleep(20 * time.Microsecond)
					docid := fmt.Sprintf("doc-%025d", i+offset+docOffset)
					err := b.Set(docid, 0, map[string]interface{}{"name": randString(sz), "junk": fmt.Sprintf("%0*d", junkSz, 0)})
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
