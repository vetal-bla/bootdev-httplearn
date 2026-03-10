package main

import (
	"fmt"
	"net/http"
)

func main() {
	servMux := http.NewServeMux()
	srv := http.Server{
		Handler: servMux,
		Addr:    ":8080",
	}
	fmt.Printf("Start server on: %s", srv.Addr)
	srv.ListenAndServe()

}
