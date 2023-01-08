package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	http.ListenAndServe("localhost:8080", mux)
}
