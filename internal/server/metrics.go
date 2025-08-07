package server

// go standard library
import (
	"fmt"
	"net/http"
)

// MiddlewareMetricsInc - increments the number of server hits for all paths /app
func (cfg *ServerConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// MetricsHandler - endpoint to show the number of hits to all endpoints followign /app
func (cfg *ServerConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf(
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>\n",
		cfg.FileServerHits.Load())
	w.Write([]byte(res))
}
