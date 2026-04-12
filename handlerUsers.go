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
		Email:          param.Email,
		HashedPassword: hashedPass,
	}

	dbUser, err := c.db.CreateUser(req.Context(), createUserParams)
	if err != nil {
		log.Printf("Error in database query: %v", err)
		respondWithError(w, http.StatusInternalServerError, "db query error")
		return
	}

	myDbuser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, myDbuser)
	log.Printf("User added to databaes: %s", param.Email)

}

func (c *apiconfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Errors get auth header: %v", err)
		respondWithError(w, http.StatusUnauthorized, "oups!")
		return
	}

	type parameters struct {
		NewPassword string `json:"password"`
		NewEmail    string `json:"email"`
	}

	userID, err := auth.ValidateJWT(accessToken, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "oups!")
		return
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err = decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get parameters")
		return
	}

	if len(param.NewEmail) < 4 {
		log.Printf("email is too short: %s", param.NewEmail)
		respondWithError(w, http.StatusNotAcceptable, "email is to short")
		return
	}

	if len(strings.Split(param.NewEmail, "@")) < 2 {
		log.Printf("area you sure that it is email: %s", param.NewEmail)
		respondWithError(w, http.StatusNotAcceptable, "not sure that it is email")
		return
	}

	hashedPassword, err := auth.HashPassword(param.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant update password")
		return
	}

	updateParams := database.UpdateEmailAndPawwordParams{
		ID:             userID,
		Email:          param.NewEmail,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := c.db.UpdateEmailAndPawword(req.Context(), updateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem with update, sorry")
		return
	}

	user := User{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)

}
