package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// CreateUserHandler - Creates a user and adds to database
func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	// parse request json
	r := createUserRequest{}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		return
	}
	if err := json.Unmarshal(reqBody, &r); err != nil {
		log.Printf("error unmarshalling request: %v", err)
		return
	}
	// add user to database
	u, err := cfg.DBQueries.CreateUser(req.Context(), r.Email)
	if err != nil {
		log.Printf("error adding user to database, %v", err)
		return
	}
	// return user object as json in response
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := createUserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	resJSON, err := json.Marshal(&res)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		return
	}
	w.Write(resJSON)
}
