package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
	"github.com/vetal-bla/bootdev-httplearn/internal/database"
)

func (c *apiconfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Cant get token form header: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error")
		return
	}

	revokeTokensParams := database.RevokeRefreshTokenParams{
		Token: refreshToken,
		UpdatedAt:  time.Now(),
		RevokedAt:  sql.NullTime{
			Time:  time.Now(),
			Valid:  true,
		},
	}

	err = c.db.RevokeRefreshToken(req.Context(), revokeTokensParams)
	if err != nil {
		log.Printf("Can't revoke token: %v", refreshToken)
		respondWithError(w, http.StatusInternalServerError, "error")
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("OK"))

}
