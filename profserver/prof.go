package profserver

import _ "net/http/pprof"
import "net/http"
import "runtime"
import "encoding/json"

func Init(Addr string) {
	http.HandleFunc("/stats/mem", handleMemStats)
	go func() {
		http.ListenAndServe(Addr, nil)
	}()
}

func handleMemStats(w http.ResponseWriter, r *http.Request) {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	bytes, _ := json.Marshal(stats)
	w.WriteHeader(200)
	w.Write(bytes)
}
