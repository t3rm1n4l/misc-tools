package profserver

import _ "net/http/pprof"
import "net/http"

var Addr = "0:9102"

func runProf() {
	go func() {
		http.ListenAndServe(Addr, nil)
	}()
}
