package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {

	var badWords = []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type parameters struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
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

	cleanedBody := replaceBadWords(badWords, param.Body)

	respondWithJSON(w, 200, responseBody{
		CleanedBody: cleanedBody,
	})
}

func replaceBadWords(badWords []string, originalString string) string {
	originalStringAr := strings.Split(originalString, " ")
	newString := make([]string, len(originalStringAr))
	for i, word := range originalStringAr {
		if slices.Contains(badWords, strings.ToLower(word)) {
			newString[i] = "****"
		} else {
			newString[i] = word
		}
	}

	return strings.Join(newString, " ")
}
