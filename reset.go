package main

import "net/http"

func (c *apiconfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	c.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
