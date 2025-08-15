package server

import (
	"log"
	"net/http"

	"github.com/k4rldoherty/go-http-server/internal/auth"
)

func (cfg *ServerConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("REVOKE - %v", err)
		return
	}

	err = cfg.DBQueries.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		log.Printf("REVOKE - %v", err)
		return
	}

	w.WriteHeader(204)
}
