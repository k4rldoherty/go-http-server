package server

import (
	"fmt"
	"log"
	"net/http"
)

// ResetHandler - resets the hits number to 0
func (cfg *ServerConfig) ResetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(403)
		return
	}
	// remove all users from db
	err := cfg.DBQueries.DeletAllUsers(req.Context())
	if err != nil {
		log.Printf("error deleting users from database: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.FileServerHits.Swap(0)
	res := fmt.Sprintf("Hits: %v", cfg.FileServerHits.Load())
	w.Write([]byte(res))
}
