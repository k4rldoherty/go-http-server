package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/auth"
	"github.com/k4rldoherty/go-http-server/internal/database"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// CreateUserHandler - Creates a user and adds to database
func (cfg *ServerConfig) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
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
	// hash password
	hash, err := auth.HashPassword(r.Password)
	if err != nil {
		log.Printf("error hashing password: %v\n", err)
		return
	}
	// add user to database
	u, err := cfg.DBQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          r.Email,
		HashedPassword: hash,
	})
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
	log.Println("user created successfully")
}
