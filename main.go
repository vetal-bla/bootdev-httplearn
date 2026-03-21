package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiconfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const serverPort = ":8080"

	config := apiconfig{
		fileServerHits: atomic.Int32{},
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
