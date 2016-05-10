package main

import "crypto/rand"
import rnd "math/rand"
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

	if len(os.Args) != 11 {
		fmt.Println("./randdocs cluster:port bucket count docid_sz field_sz is_arr junk_field_sz iterations threads doc_offset")
		os.Exit(1)
	}

	cluster := os.Args[1]
	c, _ := couchbase.Connect(fmt.Sprintf("http://%s", cluster))
	p, _ := c.GetPool("default")
	b, _ := p.GetBucket(os.Args[2])

	n, _ := strconv.Atoi(os.Args[3])
	docidSz, _ := strconv.Atoi(os.Args[4])
	sz, _ := strconv.Atoi(os.Args[5])
	is_arr := os.Args[6]
	junkSz, _ := strconv.Atoi(os.Args[7])
	nitr, _ := strconv.Atoi(os.Args[8])
	nthr, _ := strconv.Atoi(os.Args[9])
	docOffset, _ := strconv.Atoi(os.Args[10])

	runtime.GOMAXPROCS(nthr)

	for x := 0; x < nitr; x++ {
		var wg sync.WaitGroup
		for y := 0; y < nthr; y++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				for i := 0; i < n/nthr; i++ {
					time.Sleep(20 * time.Microsecond)
					docid := fmt.Sprintf("doc-%0*d", docidSz, i+offset+docOffset)
					value := make(map[string]interface{})
					value["name"] = randString(sz)
					if junkSz != 0 {
						value["junk"] = fmt.Sprintf("%0*d", junkSz, 0)
					}

					if is_arr == "yes" {
						seed := rnd.Int() % 1000000
						value["arr"] = []int{seed, seed + 1, seed + 2, seed + 3, seed + 4, seed + 5, seed + 6, seed + 7, seed + 8, seed + 9}
					}
					_ = sz
					_ = junkSz
					err := b.Set(docid, 0, value)
					//time.Sleep(time.Microsecond * 600)
					//time.Sleep(time.Microsecond * 600)
					//err := b.Delete(docid)
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
