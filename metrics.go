package main

import (
	"fmt"
	"log"
	"net/http"
)

func (c *apiconfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responceText := fmt.Sprintf("Hits: %d\n", c.fileServerHits.Load())
	w.Write([]byte(responceText))
}

func (c *apiconfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		log.Printf("%s %s. Count: %d", r.Method, r.URL.Path, c.fileServerHits.Load())
		next.ServeHTTP(w, r)
	})
}
