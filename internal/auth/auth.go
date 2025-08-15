// Package auth - all things related to authorisation including jwt implementation
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWTConfig - static properties related to the generation and validation of the jwt
type JWTConfig struct {
	SigningString []byte
	Duration      time.Duration
	Issuer        string
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, cfg *JWTConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	})
	ss, err := token.SignedString(cfg.SigningString)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString string, cfg *JWTConfig) (uuid.UUID, error) {
	// initialize an empty struce to be populated by the data read from the jwt
	claims := &jwt.RegisteredClaims{}
	// parse the token
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// optionally check signing method in here
		return cfg.SigningString, nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	s, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", errors.New("cannot find header bearer")
	}
	// Bearer TOKEN
	b := strings.Split(bearer, " ")
	if len(b) != 2 {
		return "", errors.New("cannot find header bearer")
	}
	return b[1], nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
