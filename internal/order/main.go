package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	log.Println("Listenning on 8082...")
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println("trying to pong")
		_, _ = io.WriteString(w, "pong")
	})
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
