package handlers

import (
	"log"
	"strconv"
	"encoding/json"
	"net/http"
	"roles/cache"
	"roles/db"
	"roles/models"

    "github.com/go-chi/chi/v5"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	cached := cache.GetRoles()
	json.NewEncoder(w).Encode(cached)
}

func GetRoleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {		
		log.Println("Bad Id")
		http.NotFound(w, r)
		return 
	}
	role, err := cache.FindById(uint(id))
	if err != nil {
		log.Println("Role Not foun")
		http.NotFound(w, r)
        return
	}

	json.NewEncoder(w).Encode(role)
}

//tools

func LoadCache() ([]models.Role) {
	var roles []models.Role
	if err := db.DB.Find(&roles).Error; err != nil {
		log.Println("db empty")
		return roles
	}

	cache.SetRoles(roles)
	return roles
}
