package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
)

func (c *apiconfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
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

	hashedPass, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("Cant create hash: %v", err)
		respondWithError(w, http.StatusInternalServerError, "sorry")
		return
	}

	createUserParams := database.CreateUserParams{
		Email:  param.Email,
		HashedPassword:  hashedPass,
	}

	dbUser, err := c.db.CreateUser(req.Context(), createUserParams)
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
