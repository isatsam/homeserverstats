package node

import (
	"encoding/json"
	"net/http"
)

func StartNode() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getStats)
	http.ListenAndServe("localhost:8001", mux)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	s := newStats("node0")
	j, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	w.Write(j)
}
