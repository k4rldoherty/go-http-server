package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	res := loginResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
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
