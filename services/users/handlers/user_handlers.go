package handlers

import (
    "encoding/json"
    "net/http"
    "users/cache"
    "users/db"
    "users/models"
	"log"
	"time"

    "github.com/go-chi/chi/v5"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all Users")
    var users []models.User
    db.DB.Find(&users)
    json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
	cacheKey := "user:" + id
    var user models.User
    // Buscar en caché
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
        json.Unmarshal([]byte(cached), &user)
        json.NewEncoder(w).Encode(user)
		return
	}

   
    if err := db.DB.First(&user, id).Error; err != nil {
        http.NotFound(w, r)
        return
    }

    // Serializar y guardar en caché
	data, _ := json.Marshal(user)
	cache.Set(cacheKey, string(data), 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)
    db.DB.Create(&user)
    json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    db.DB.Delete(&models.User{}, id)
    cache.Delete("user:" + id)
    w.WriteHeader(http.StatusNoContent)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var user models.User
    if err := db.DB.First(&user, id).Error; err != nil {
        http.NotFound(w, r)
        return
    }

    var updated models.User
    json.NewDecoder(r.Body).Decode(&updated)

    user.Username = updated.Username
    user.Password = updated.Password
    user.RoleID = updated.RoleID

    db.DB.Save(&user)
    cache.Delete("user:" + id)

    json.NewEncoder(w).Encode(user)
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "username query param is required", http.StatusBadRequest)
        return
    }

    var users []models.User
    db.DB.Where("username ILIKE ?", "%"+username+"%").Find(&users)
    json.NewEncoder(w).Encode(users)
}
