package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/products", GetProductsHandler)
	r.Get("/products/{id}", PrintSomething)

	log.Println("Product service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
