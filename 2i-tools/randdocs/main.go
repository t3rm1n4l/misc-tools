package main

import "crypto/rand"
import "github.com/couchbase/go-couchbase"
import "os"
import "strconv"

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

	c, _ := couchbase.Connect("http://Administrator:asdasd@127.0.0.1:9000")
	p, _ := c.GetPool("default")
	b, _ := p.GetBucket("default")

	n, _ := strconv.Atoi(os.Args[1])
	sz, _ := strconv.Atoi(os.Args[2])

	for i := 0; i < n; i++ {
		docid := randString(sz)
		b.Set(docid, 0, map[string]interface{}{"name": docid})
	}

}
