package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	RoleName     string `json:"role,omitempty"`
	PasswordHash string `json:"-"`
}

// Handler para registrar usuario
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if u.Username == "" || u.Password == "" || u.RoleName == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	err := registerUser(u.Username, u.Password, u.RoleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Handler para login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	username, role, err := authenticateUser(u.Username, u.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(username)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"username": username,
		"role":     role,
		"token":    token,
	}
	json.NewEncoder(w).Encode(resp)
}

// Funciones para registrar y autenticar

func registerUser(username, password, role string) error {
	var roleID int
	err := db.QueryRow("SELECT id FROM roles WHERE name=$1", role).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash, role_id) VALUES ($1, $2, $3)", username, string(hashed), roleID)
	return err
}

func authenticateUser(username, password string) (string, string, error) {
	var storedHash, roleName string
	err := db.QueryRow(`
        SELECT u.password_hash, r.name 
        FROM users u JOIN roles r ON u.role_id = r.id 
        WHERE u.username = $1`, username).Scan(&storedHash, &roleName)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid password")
	}
	return username, roleName, nil
}
