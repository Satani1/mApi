package main

import (
	"api/api"
	"net/http"
)

func main() {
	srv := api.NewServer()
	http.ListenAndServe("localhost:8080", srv)
	//ruslan
}
