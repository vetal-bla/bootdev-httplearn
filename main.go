package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type apiconfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const serverPort = ":8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("cant connect to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	config := apiconfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	servMux := http.NewServeMux()
	fileSrv := config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	servMux.Handle("/app/", fileSrv)

	servMux.HandleFunc("GET /api/healthz", handlerHealthz)
	servMux.HandleFunc("POST /api/users", config.handlerCreateUser)
	servMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	servMux.HandleFunc("GET /admin/metrics", config.handlerMetrics)
	servMux.HandleFunc("POST /admin/reset", config.handlerReset)

	srv := http.Server{
		Handler: servMux,
		Addr:    serverPort,
	}

	fmt.Printf("Start server on port: %s\nWhich server dir: %s\n", serverPort, filepathRoot)
	log.Fatal(srv.ListenAndServe())

}
