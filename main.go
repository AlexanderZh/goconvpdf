package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("Server started")

	router := NewRouter()

	serve, ok := os.LookupEnv("PORT")
	if !ok {
		serve = "8080"
	}
	serve = ":" + serve
	log.Fatal(http.ListenAndServe(serve, router))

}
