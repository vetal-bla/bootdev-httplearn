package main

import (
	"fmt"
	"net/http"
)

func main() {
	filepathRoot := "."
	serverPort := ":8080"

	servMux := http.NewServeMux()
	servMux.Handle("/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	servMux.HandleFunc("/healthz", handlerHealthz)
	srv := http.Server{
		Handler: servMux,
		Addr:    serverPort,
	}
	fmt.Printf("Start server on port: %s\nWhich server dir: %s\n", serverPort, filepathRoot)
	srv.ListenAndServe()

}

func handlerHealthz(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("test - handler")
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}
