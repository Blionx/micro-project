package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
    "net/url"
    "strings"

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

	r.Mount("/users", UsersRouter())

	return r
}

func UsersRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", proxyHandler("http://users:8083"))
	r.Get("/search", proxyHandler("http://users:8083"))
	r.Get("/{id}", proxyHandler("http://users:8083"))
	r.Post("/", proxyHandler("http://users:8083"))
	r.Put("/{id}",  proxyHandler("http://users:8083"))
	r.Delete("/{id}", proxyHandler("http://users:8083"))
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

func proxyHandler(targetBase string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        targetURL, err := url.Parse(targetBase)
        if err != nil {
            http.Error(w, "Bad proxy target", http.StatusInternalServerError)
            return
        }

        // Crear una nueva solicitud con misma ruta
        proxy := httputil.NewSingleHostReverseProxy(targetURL)

        // Esto asegura que la ruta original (con {id}, etc.) se respete
        originalDirector := proxy.Director
        proxy.Director = func(req *http.Request) {
            originalDirector(req)
            req.URL.Path = singleJoiningSlash(targetURL.Path, r.URL.Path)
            req.URL.RawQuery = r.URL.RawQuery
            req.Host = targetURL.Host
        }

        proxy.ServeHTTP(w, r)
    }
}

// helper para evitar dobles slashes
func singleJoiningSlash(a, b string) string {
    aslash := strings.HasSuffix(a, "/")
    bslash := strings.HasPrefix(b, "/")
    switch {
    case aslash && bslash:
        return a + b[1:]
    case !aslash && !bslash:
        return a + "/" + b
    }
    return a + b
}
