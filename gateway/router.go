package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Rutas públicas (sin autenticación)
	r.Post("/auth/login", ProxyToAuth)
	r.Post("/user/register", ProxyToNewUser)

	// Rutas protegidas con JWT
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware) // Middleware de autenticación para estas rutas

		r.Get("/products", ProxyToProducts)
		r.Get("/products/{id}", ProxyToProducts)
		// Puedes agregar más rutas protegidas aquí
	})

	return r
}

func ProxyToNewUser(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post("http://auth:8081/register", "application/json", r.Body)

	if err != nil {
		http.Error(w, "Error Registering user", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func ProxyToAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("routing request to login.")
	resp, err := http.Post("http://auth:8081/login", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Auth service error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func ProxyToProducts(w http.ResponseWriter, r *http.Request) {
	// Cambia según tu servicio real
	req, err := http.NewRequest(r.Method, "http://products:8082"+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Request error", http.StatusInternalServerError)
		return
	}

	// Copia headers
	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Products service error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
