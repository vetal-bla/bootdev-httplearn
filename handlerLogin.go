package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
)

func (c *apiconfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get parameters")
		return
	}

	dbUser, err := c.db.GetUserByMail(req.Context(), param.Email)
	if err != nil {
		log.Printf("Error retrive user by email: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	passwordOK, err := auth.CheckPasswordHash(param.Password, dbUser.HashedPassword)

	if err != nil {
		log.Printf("Cant check password: %v", err)
		return
	}

	if !passwordOK {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	jsonUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	log.Printf("Hash matched. OK!")
	respondWithJSON(w, http.StatusOK, jsonUser)
}
