package main

import (
	"log"
	"net/http"
)

const proxyAddr string = "localhost:9000"

var (
	counter    int    = 0
	FirstAddr  string = "localhost:9001"
	SecondAddr string = "localhost:9002"
)

func main() {

}

func HandleProxy(w http.ResponseWriter, r *http.Request) {

	if counter == 0 {
		_, err := http.Post(FirstAddr, "application/json", r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		counter++
		return
	}

	counter--
}
