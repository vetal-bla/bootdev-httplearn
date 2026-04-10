package main

import (
	"log"
	"net/http"
	"time"

	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
)

func (c *apiconfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {	
	type responce struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Cant get token form header: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error")
		return
	}

	refreshTokenDB, err := c.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Database query error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "database error")
		return
	}

	if refreshTokenDB.Token == "" {
		log.Printf("Empty token:%v", refreshTokenDB.Token)
		respondWithError(w, http.StatusUnauthorized, "database error")
		return
	}
	if refreshTokenDB.RevokedAt.Valid {
		log.Printf("Token revoked at: %v", refreshTokenDB.RevokedAt)
		respondWithError(w, http.StatusUnauthorized, "database error")
		return
	}
	if refreshTokenDB.ExpiresAt.Before(time.Now()){
		log.Printf("Token is expired: %v", refreshTokenDB.ExpiresAt)
		respondWithError(w, http.StatusUnauthorized, "database error")
		return
	}

	token, err := auth.MakeJWT(refreshTokenDB.UserID, c.secret, time.Hour)
	if err != nil {
		log.Printf("Cant get token form header: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error")
		return
	}

	log.Printf("token: %s", token)

	t := responce{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, t)

}
