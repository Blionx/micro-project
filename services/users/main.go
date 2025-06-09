package main

import (
    "log"
    "net/http"
    "os"
    "users/db"
    "users/handlers"
    "users/cache"

    "github.com/go-chi/chi/v5"
)

func main() {
    if err := db.ConnectDB(); err != nil {
        log.Fatalf("DB error: %v", err)
    }
    cache.InitRedis()
	
    r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
    r.Get("/", handlers.GetUsers)          // <-- esta debe estar primero
    r.Get("/{id}", handlers.GetUserByID)   // <-- esta después
	r.Get("/search", handlers.SearchUser)
    r.Post("/", handlers.CreateUser)
    r.Put("/{id}", handlers.UpdateUser)
    r.Delete("/{id}", handlers.DeleteUser)
})


    port := os.Getenv("PORT")
    if port == "" {
        port = "8083"
    }

    log.Printf("User service running on port %s", port)
    http.ListenAndServe(":"+port, r)
}
