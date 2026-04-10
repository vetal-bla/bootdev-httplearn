package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
)

func (c *apiconfig) handlerCreateChirps(w http.ResponseWriter, req *http.Request) {

	var badWords = []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type parameters struct {
		Body   string    `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err !=nil {
		respondWithError(w, http.StatusInternalServerError, "token error")
		log.Printf("Get token erro: %v", err)
		return
	}

	userID, err := auth.ValidateJWT(token, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		log.Printf("JWT validation error: %v\n", err)
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err = decoder.Decode(&param)
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

	dbParams := database.CreateChirpsParams{
		Body:   cleanedBody,
		UserID: userID,
	}

	dbChirps, err := c.db.CreateChirps(req.Context(), dbParams)
	if err != nil {
		log.Printf("Error in database query: %v", err)
		respondWithError(w, http.StatusInternalServerError, "db query error")
	}

	respondWithJSON(w, http.StatusCreated, Chirps{
		ID:        dbChirps.ID,
		CreatedAt: dbChirps.CreatedAt,
		UpdatedAt: dbChirps.UpdatedAt,
		Body:      dbChirps.Body,
		UserID:    dbChirps.UserID,
	})
}

func (c *apiconfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := c.db.GetAllChirps(req.Context())
	if err != nil {
		log.Printf("Error retrieve data from database: %v", err)
		return
	}

	arrChirps := []Chirps{}
	for _, v := range dbChirps {
		arrChirps = append(arrChirps, Chirps{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserID:    v.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, arrChirps)
}

func (c *apiconfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {

	chirpID := req.PathValue("chirpid")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("Cant convert string to uuid: %v", err)
		respondWithError(w, http.StatusNotAcceptable, "Uncorrect uuid format")
		return
	}
	fmt.Println(chirpID)

	dbChirp, err := c.db.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		log.Printf("Cant retrieve data from database: %v", err)
		respondWithError(w, http.StatusNotFound, "Not found!")
		return
	}
	fmt.Println(dbChirp)

	respondWithJSON(w, http.StatusOK, Chirps{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
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
