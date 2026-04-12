package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (c *apiconfig) handlerWebhook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get parameters")
		return
	}

	if param.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "event not user.upgraded")
		return
	}

	paramUserID, err := uuid.Parse(param.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "uuid wrong format")
		return
	}

	user, err := c.db.UpdateChirpyRed(req.Context(), paramUserID)
	if err != nil {
		log.Printf("Error update user chirpy red: %v", err)
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}

	if user.ID == paramUserID {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

}
