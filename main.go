package main

// go standard library
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	fileServerHandler := http.FileServer((http.Dir(filePathRoot)))
	// a serveMux directs traffic to the relevant handler function - like a controller in .NET
	serveMux := http.NewServeMux()
	// sets /app as the root directory, even when you dont actually have an app directory on your server
	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServerHandler)),
	)

	// util endpoints to check health etc
	serveMux.HandleFunc("GET /api/healthz", healthzHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", cfg.validate_chirpHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	// server obejct to listen on port 8080
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %v and listening on port %v", filePathRoot, port)
	// starts the server
	log.Fatal(server.ListenAndServe())
}

// increments the number of server hits for all paths /app
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// endpoint to show the number of hits to all endpoints followign /app
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>\n", cfg.fileServerHits.Load())
	w.Write([]byte(res))
}

func (cfg *apiConfig) validate_chirpHandler(w http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}
	type validReturn struct {
		Body bool `json:"valid"`
	}
	type invalidReturn struct {
		Body string `json:"error"`
	}
	decoder := json.NewDecoder(req.Body)
	rb := reqBody{}
	err := decoder.Decode(&rb)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		respBody := invalidReturn{Body: "Something went wrong"}
		dat, _ := json.Marshal(respBody)
		w.Write([]byte(dat))
		return
	}

	if (len(rb.Body)) > 140 {
		w.WriteHeader(400)
		respBody := invalidReturn{Body: "Chirp is too long"}
		dat, _ := json.Marshal(respBody)
		w.Write([]byte(dat))
		return
	}

	w.WriteHeader(200)
	respBody := validReturn{Body: true}
	dat, _ := json.Marshal(respBody)
	w.Write([]byte(dat))
}
