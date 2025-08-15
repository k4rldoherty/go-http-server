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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *ServerConfig) LoginHandler(r http.ResponseWriter, req *http.Request) {
	// parse Request
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("LOGIN - %v", err)
		return
	}
	loginReq := loginRequest{}
	if err = json.Unmarshal(reqBody, &loginReq); err != nil {
		log.Printf("LOGIN - %v", err)
		return
	}
	// look up user by email
	u, err := cfg.DBQueries.GetUserByEmail(req.Context(), loginReq.Email)
	if err != nil {
		log.Printf("LOGIN - %v", err)
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
		log.Printf("LOGIN - %v", err)
		return
	}
	// generate refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("LOGIN - %v", err)
		return
	}

	// add refresh token to DB
	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // 60 days
	}
	_, err = cfg.DBQueries.CreateRefreshToken(req.Context(), params)
	if err != nil {
		log.Printf("LOGIN - %v", err)
		return
	}

	// would save the token to db around here
	res := loginResponse{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  u.IsChirpyRed,
	}

	// format the response as json
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("LOGIN - %v", err)
		return
	}
	r.Header().Add("Content-Type", "application/json")
	r.WriteHeader(200)
	r.Write(resJSON)
}
