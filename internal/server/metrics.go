package server

// go standard library
import (
	"fmt"
	"net/http"
)

// increments the number of server hits for all paths /app
func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// endpoint to show the number of hits to all endpoints followign /app
func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf(
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>\n",
		cfg.FileServerHits.Load())
	w.Write([]byte(res))
}
