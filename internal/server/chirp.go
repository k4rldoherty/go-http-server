package server

// go standard library
import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/auth"
	"github.com/k4rldoherty/go-http-server/internal/database"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

type reqBody struct {
	Body string `json:"body"`
}

type resBody struct {
	Body   string    `json:"body"`
	UserID string    `json:"user_id"`
	ID     uuid.UUID `json:"id"`
}

type errBody struct {
	Body string `json:"body"`
}

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ServerConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) {
	// authenticate
	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("CHIRP GET BEARER - %v", err)
		w.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(bearer, cfg.JWTCfg)
	if err != nil {
		log.Printf("CHIRP VALIDATE JWT - %v", err)
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(req.Body)
	rb := reqBody{}
	err = decoder.Decode(&rb)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Printf("CHIRP - %v", err)
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
		UserID: userID,
	})
	if err != nil {
		log.Printf("CHIRP - %v", err)
		w.WriteHeader(400)
		return
	}
	respBody := resBody{Body: c.Body, UserID: userID.String(), ID: c.ID}
	dat, _ := json.Marshal(respBody)
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *ServerConfig) GetAllChirpsHandler(w http.ResponseWriter, req *http.Request) {
	c, err := cfg.DBQueries.GetAllChirps(req.Context())
	if err != nil {
		w.WriteHeader(400)
		log.Printf("CHIRP - %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	res := []chirp{}
	for _, v := range c {
		parsedItem := chirp{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserID:    v.UserID,
		}
		res = append(res, parsedItem)
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(400)
		log.Printf("CHIRP - %v", err)
		return
	}
	w.Write(resJSON)
}

func (cfg *ServerConfig) GetChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	if id == "" {
		log.Println("CHIRP - id passed in was empty")
		return
	}
	log.Printf("should be a uuid -> %v\n", id)
	idParsed, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(400)
		log.Printf("CHIRP - %v", err)
		return
	}
	// try get chirp by id
	c, err := cfg.DBQueries.GetChirpById(req.Context(), idParsed)
	if err != nil {
		log.Printf("CHIRP - %v", err)
		w.WriteHeader(404)
		return
	}
	parsedChirp, err := json.Marshal(chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	})
	if err != nil {
		w.WriteHeader(400)
		log.Printf("CHIRP - %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	w.Write(parsedChirp)
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
