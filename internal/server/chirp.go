package server

// go standard library
import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/database"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

type reqBody struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type resBody struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type errBody struct {
	Body string `json:"body"`
}

func (cfg *APIConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	rb := reqBody{}
	err := decoder.Decode(&rb)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		respBody := reqBody{Body: "Something went wrong"}
		dat, _ := json.Marshal(respBody)
		w.Write([]byte(dat))
		return
	}

	if (len(rb.Body)) > 140 {
		w.WriteHeader(400)
		respBody := errBody{Body: "Chirp is too long"}
		dat, _ := json.Marshal(respBody)
		w.Write([]byte(dat))
		return
	}

	cleanedResult := cleanInput(rb.Body)
	// add chirp to database
	c, err := cfg.DBQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedResult,
		UserID: cfg.UserID,
	})
	if err != nil {
		log.Printf("failed to create chirp: %v", err)
		w.WriteHeader(400)
		return
	}
	respBody := resBody{Body: c.Body, UserID: c.UserID}
	dat, _ := json.Marshal(respBody)
	w.WriteHeader(201)
	w.Write(dat)
	log.Println("chirp created successfully.")
}

func cleanInput(s string) string {
	cleanedResult := make([]string, 0)
	for word := range strings.SplitSeq(s, " ") {
		wordLowered := strings.ToLower(word)
		if slices.Contains(profaneWords, wordLowered) {
			cleanedResult = append(cleanedResult, "****")
		} else {
			cleanedResult = append(cleanedResult, word)
		}
	}
	return strings.Join(cleanedResult, " ")
}
