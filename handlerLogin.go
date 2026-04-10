package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
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

	jwtToken, err := auth.MakeJWT(dbUser.ID, c.secret, time.Hour)
	if err != nil {
		log.Printf("Cant create jwt token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Problem with auth")
		return
	}

	dbRefreshTokensParams := database.CreateRefreshTokensParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	dbRefreshToken, err := c.db.CreateRefreshTokens(req.Context(), dbRefreshTokensParams)
	if err != nil {
		log.Printf("Cant insert refresh token to database: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Problem with auth")
		return
	}

	fmt.Printf("database created token: %v", dbRefreshToken)

	jsonUser := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        jwtToken,
		RefreshToken: dbRefreshToken.Token,
	}

	log.Printf("Hash matched. OK! Token created for user: %s", jsonUser.ID)
	respondWithJSON(w, http.StatusOK, jsonUser)
}
