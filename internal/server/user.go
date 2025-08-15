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

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type editUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserHandler - Creates a user and adds to database
func (cfg *ServerConfig) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	// parse request json
	r := createUserRequest{}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("CREATE USER - %v", err)
		return
	}
	if err := json.Unmarshal(reqBody, &r); err != nil {
		log.Printf("CREATE USER - %v", err)
		return
	}
	// hash password
	hash, err := auth.HashPassword(r.Password)
	if err != nil {
		log.Printf("CREATE USER - %v\n", err)
		return
	}
	// add user to database
	u, err := cfg.DBQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          r.Email,
		HashedPassword: hash,
	})
	if err != nil {
		log.Printf("CREATE USER - %v", err)
		return
	}
	// return user object as json in response
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := userResponse{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
	resJSON, err := json.Marshal(&res)
	if err != nil {
		log.Printf("CREATE USER - %v", err)
		return
	}
	w.Write(resJSON)
}

func (cfg *ServerConfig) EditUserHandler(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	parsedReq := editUserRequest{}
	err = json.Unmarshal(reqBody, &parsedReq)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.JWTCfg)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	hashedPwd, err := auth.HashPassword(parsedReq.Password)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	params := database.EditUserParams{
		ID:             userID,
		Email:          parsedReq.Email,
		HashedPassword: hashedPwd,
	}
	updatedDetails, err := cfg.DBQueries.EditUser(req.Context(), params)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	res := userResponse{
		ID:          userID,
		CreatedAt:   updatedDetails.CreatedAt,
		UpdatedAt:   updatedDetails.UpdatedAt,
		Email:       updatedDetails.Email,
		IsChirpyRed: updatedDetails.IsChirpyRed,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("EDIT USER - %v", err)
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(200)
	w.Write(resBytes)
}
