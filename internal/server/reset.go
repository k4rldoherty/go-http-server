package server

import (
	"fmt"
	"net/http"
)

// resets the hits number to 0
func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.FileServerHits.Swap(0)
	res := fmt.Sprintf("Hits: %v", cfg.FileServerHits.Load())
	w.Write([]byte(res))
}
