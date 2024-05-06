package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

	log.Println("running server at :3000")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
