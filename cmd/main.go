package main

import (
	"log"
	"net/http"
	"os"

	"dms/internal"

	"github.com/joho/godotenv"
)

var store *internal.Store
var jwtSecret []byte

func main() {
	// Load environment variables
	godotenv.Load()
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	pgDSN := os.Getenv("PG_DSN")
	redisAddr := os.Getenv("REDIS_ADDR")

	// Initialize Store
	var err error
	store, err = internal.NewStore(pgDSN, redisAddr)
	if err != nil {
		log.Fatal("Cannot connect to store:", err)
	}
	log.Println("Connected to Postgres + Redis!")

	// Router
	r := internal.NewRouter(store, jwtSecret)

	// Start HTTP server
	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}