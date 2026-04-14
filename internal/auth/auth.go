package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("Password can't be empty")
	}

	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Cant craete hash from password: %v", err)
	}

	return hashedPass, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Cant compare password with hash: %v", err)
	}

	return match, nil
}

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)

	if err != nil || !token.Valid {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("authorization")
	if authHeader == "" {
		return "", errors.New("Empty auth header")
	}

	token := strings.Fields(authHeader)
	if len(token) <= 1 {
		return "", errors.New("Uncorrect token")
	}
	return token[1], nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("authorization")
	if authHeader == "" {
		return "", errors.New("Empty auth header")
	}

	token := strings.Fields(authHeader)
	if len(token) < 2 {
		return "", errors.New("Uncorrect token format")
	}

	if token[0] != "ApiKey" {
		return "", errors.New("Uncorrect token format")
	}

	return token[1], nil
}

func MakeRefreshToken() string {
	refreshToken := make([]byte, 32)
	rand.Read(refreshToken)
	return hex.EncodeToString(refreshToken)
}
