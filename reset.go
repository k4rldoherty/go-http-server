package main

import (
	"fmt"
	"net/http"
)

// resets the hits number to 0
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Swap(0)
	res := fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())
	w.Write([]byte(res))
}
