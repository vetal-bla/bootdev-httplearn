package main

import (
	"fmt"
	"log"
	"net/http"
)

func (c *apiconfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Printf("Running platform: %s", c.platform)

	if c.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	c.fileServerHits.Store(0)
	err := c.db.DeleteAllUsers(req.Context())
	if err != nil {
		log.Printf("Cant remove users from database:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
