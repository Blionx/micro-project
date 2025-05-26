package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/verify", VerifyHandler)

	log.Println("Auth service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
