package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Connected to the database.")
				return db, nil
			}
		}
		log.Printf("Waiting for the database... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after 10 attempts: %w", err)
}

func initRoles() error {
	roles := []string{"admin", "user"}

	for _, role := range roles {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", role).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err := db.Exec("INSERT INTO roles (name) VALUES ($1)", role)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/register", RegisterHandler)
	r.Post("/login", LoginHandler)

	log.Println("Auth service running on :8081")
	http.ListenAndServe(":8081", r)
}
