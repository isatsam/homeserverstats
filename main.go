package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	fmt.Println("serving on http://localhost:8000")
	http.ListenAndServe("localhost:8000", mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}
