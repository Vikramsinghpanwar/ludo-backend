package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Ludo backend...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
