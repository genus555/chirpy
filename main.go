package main

import _ "github.com/lib/pq"

import (
	//"fmt"
	"github.com/joho/godotenv"
	"github.com/genus555/chirpy/internal/database"
	
	"sync/atomic"
	"os"
	"log"
	"database/sql"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db				*database.Queries
	ts				string
}
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("tokenSecret")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print(err)
	}
	dbQueries := database.New(db)

	var serverHits atomic.Int32
	initializeHits(&serverHits)
	cfg := apiConfig{
		fileserverHits: serverHits,
		db:				dbQueries,
		ts:	tokenSecret,
	}
	launchServer(&cfg);
}

func initializeHits(hits *atomic.Int32) {
	hits.Store(int32(0))
}