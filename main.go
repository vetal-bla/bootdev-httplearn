package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
)

type apiconfig struct {
	fileServerHits atomic.Int32
	db	*database.Queries
}

func main() {
	const filepathRoot = "."
	const serverPort = ":8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
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
		db: dbQueries,
	}


	servMux := http.NewServeMux()
	fileSrv := config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	servMux.Handle("/app/", fileSrv)

	servMux.HandleFunc("GET /api/healthz", handlerHealthz)
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
