package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get parameters")
		return
	}

	const maxChirpLength = 140

	if len(param.Body) > maxChirpLength {
		log.Printf("Chirp is too long")
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, responseBody{
		Valid: true,
	})
}
