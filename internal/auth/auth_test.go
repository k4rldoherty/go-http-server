package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TODO:
// Add tests for hash password
// Add tests for check password hash

func TestValidateJWT(t *testing.T) {
	testJWTCfg := &JWTConfig{
		Issuer:        "chirpy",
		Duration:      time.Second + 5,
		SigningString: []byte("secret"),
	}
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, testJWTCfg)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret []byte
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: testJWTCfg.SigningString,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: testJWTCfg.SigningString,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: []byte("wrong_secret"),
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, testJWTCfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if tt.name == "Invalid token" {
				testJWTCfg.SigningString = []byte("invalid")
			}
		})
	}
}

func GetBearerTokenTest(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		header  http.Header
	}{
		{
			name:    "valid header",
			wantErr: false,
		},
		{
			name:    "invalid header",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// TODO:
		// Create header data for tests
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken error = %v", err)
			}
		})
	}
}
