package main

import (
	"fmt"
	"net/http"
)

func main() {
	filepathRoot := "."
	serverPort := ":8080"

	servMux := http.NewServeMux()
	servMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	srv := http.Server{
		Handler: servMux,
		Addr:    serverPort,
	}
	fmt.Printf("Start server on port: %s\nWhich server dir: %s\n", serverPort, filepathRoot)
	srv.ListenAndServe()

}
