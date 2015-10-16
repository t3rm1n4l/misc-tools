package profserver

import _ "net/http/pprof"
import "net/http"

func Init(Addr string) {
	go func() {
		http.ListenAndServe(Addr, nil)
	}()
}
