package main

import (
	"log"
	"net/http"
	"roles/db"
	"roles/handlers"
	"roles/routes"
)

func main() {
	db.ConnectDB()
	r := routes.InitRoutes()
	
	log.Println("✅ Roles service running on port 8083")
	http.ListenAndServe(":8083", r)
	_ = handlers.LoadCache()
}
