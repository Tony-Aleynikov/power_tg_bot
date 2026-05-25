package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /handler", handler)

	s := &http.Server{
		Addr: ":8080",
	}

	fmt.Println("Listening on port 8080")
	log.Fatal(s.ListenAndServe()) // почему именно log.Fatal
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Err")
		return
	}
	defer r.Body.Close()
	fmt.Println(string(body))
}
