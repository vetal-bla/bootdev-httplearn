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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirps struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type apiconfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	const filepathRoot = "."
	const serverPort = ":8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	if secret == "" {
		log.Fatal("SECRET variable must be set")
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

	servMux.HandleFunc("GET /admin/metrics", config.handlerMetrics)
	servMux.HandleFunc("GET /api/chirps", config.handlerGetAllChirps)
	servMux.HandleFunc("GET /api/chirps/{chirpid}/", config.handlerGetChirp)
	servMux.HandleFunc("GET /api/healthz", handlerHealthz)
	servMux.HandleFunc("POST /admin/reset", config.handlerReset)
	servMux.HandleFunc("POST /api/chirps", config.handlerCreateChirps)
	servMux.HandleFunc("POST /api/login", config.handlerLogin)
	servMux.HandleFunc("POST /api/refresh", config.handlerRefresh)
	servMux.HandleFunc("POST /api/revoke", config.handlerRevoke)
	servMux.HandleFunc("POST /api/users", config.handlerCreateUser)
	servMux.HandleFunc("PUT /api/users", config.handlerUpdateUser)
	servMux.HandleFunc("DELETE /api/chirps/{chirpid}", config.handlerDeleteChirp)
	servMux.HandleFunc("POST /api/polka/webhooks", config.handlerWebhook)

	srv := http.Server{
		Handler: servMux,
		Addr:    serverPort,
	}

	fmt.Printf("Start server on port: %s\nWhich server dir: %s\n", serverPort, filepathRoot)
	log.Fatal(srv.ListenAndServe())

}
