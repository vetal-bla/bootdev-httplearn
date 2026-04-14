package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vetal-bla/bootdev-httplearn/internal/auth"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := auth.HashPassword(password1)
	hash2, _ := auth.HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := auth.CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := auth.MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	// define headers for test
	validHeader := make(http.Header)
	validHeader.Set("authorization", "Bearer VALID_TOKEN$1")

	badFormat := make(http.Header)
	badFormat.Set("authorization", "Bearer   ")

	badFormatSpaces := make(http.Header)
	badFormatSpaces.Set("authorization", "Bearer    VALID_TOKEN$1  ")

	badFormatEmptyString := make(http.Header)
	badFormatEmptyString.Set("authorization", "")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "ValidHeader",
			headers: validHeader,
			want: "VALID_TOKEN$1",
			wantErr:  false,
		},
		{
			name: "EmptyHeader",
			headers: http.Header{},
			want: "",
			wantErr:  true,
		},
		{
			name: "NoToken",
			headers: badFormat,
			want: "",
			wantErr:  true,
		},
		{
			name: "AdditionalSpaces",
			headers: badFormatSpaces,
			want: "VALID_TOKEN$1",
			wantErr:  false,
		},
		{
			name: "EmptyAthorizationHeader",
			headers: badFormatEmptyString,
			want: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestGetApiKey(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "ValidHeader",
			headers: http.Header{
				"Authorization": []string{"ApiKey VALID_TOKEN$1"},
			},
			want: "VALID_TOKEN$1",
			wantErr:  false,
		},
		{
			name: "EmptyHeader",
			headers: http.Header{},
			want: "",
			wantErr:  true,
		},
		{
			name: "NoToken",
			headers: http.Header{
				"Authorization": []string{"ApiKey   "},
			},
			want: "",
			wantErr:  true,
		},
		{
			name: "AdditionalSpaces",
			headers: http.Header{
				"Authorization": []string{"ApiKey     VALID_TOKEN$1"},
			},
			want: "VALID_TOKEN$1",
			wantErr:  false,
		},
		{
			name: "EmptyAthorizationHeader",
			headers: http.Header{
				"Authorization": []string{""},
			},
			want: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetApiKey(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetApiKey() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetApiKey() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetApiKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

