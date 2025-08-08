package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/auth"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds string `json:"expires_in_seconds"`
}

type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *ServerConfig) LoginHandler(r http.ResponseWriter, req *http.Request) {
	// parse Request
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		return
	}
	loginReq := loginRequest{}
	if err = json.Unmarshal(reqBody, &loginReq); err != nil {
		log.Printf("error parsing request: %v", err)
		return
	}

	if loginReq.ExpiresInSeconds != "" {
		parsedExpiresInSeconds, err := strconv.Atoi(loginReq.ExpiresInSeconds)
		if err != nil {
			cfg.JWTCfg.Duration = time.Hour
			log.Printf("parsing failed, defaulting to one hour: %v", err)
		}
		// set the expiry to the passed in one, otherwise leave at 1 hour
		if parsedExpiresInSeconds > int(time.Second)*60*60 {
			cfg.JWTCfg.Duration = time.Hour
		}
		cfg.JWTCfg.Duration = time.Duration(parsedExpiresInSeconds)
	}
	// look up user by email
	u, err := cfg.DBQueries.GetUserByEmail(req.Context(), loginReq.Email)
	if err != nil {
		log.Printf("user not found: %v", err)
		return
	}
	// check credentials
	if err = auth.CheckPasswordHash(loginReq.Password, u.HashedPassword); err != nil {
		log.Printf("passwords do not match: %v", err)
		r.WriteHeader(401)
		r.Write([]byte("incorrect email or password."))
		return
	}
	// generate token
	token, err := auth.MakeJWT(u.ID, cfg.JWTCfg)
	if err != nil {
		log.Printf("error generating jwt token: %v", err)
		return
	}
	// would save the token to db around here
	res := loginResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
		Token:     token,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("error parsing json response: %v", err)
		return
	}
	r.Header().Add("Content-Type", "application/json")
	r.WriteHeader(200)
	r.Write(resJSON)
}
