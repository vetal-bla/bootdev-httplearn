package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (c *apiconfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get parameters")
		return
	}

	if len(param.Email) < 4 {
		log.Printf("email is too short: %s", param.Email)
		respondWithError(w, http.StatusNotAcceptable, "email is to short")
		return
	}

	if len(strings.Split(param.Email, "@")) < 2 {
		log.Printf("area you sure that it is email: %s", param.Email)
		respondWithError(w, http.StatusNotAcceptable, "not sure that it is email")
		return
	}

	dbUser, err := c.db.CreateUser(req.Context(), param.Email)
	if err != nil {
		log.Printf("Error in database query: %v", err)
		respondWithError(w, http.StatusInternalServerError, "db query error")
		return
	}

	myDbuser := User{
		ID: dbUser.ID,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
		Email:  dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, myDbuser)
	log.Printf("User added to databaes: %s", param.Email)

}
