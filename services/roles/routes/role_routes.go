package routes

import (
	"roles/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes() (chi.Router){
	r := chi.NewRouter()
	r.Route("/roles", func(r chi.Router) {
		r.Get("/", handlers.GetRoles)
		r.Get("/{id}", handlers.GetRoleById)
	})
	return r
}