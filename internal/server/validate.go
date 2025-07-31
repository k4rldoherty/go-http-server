package server

// go standard library
import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

type reqBody struct {
	Body string `json:"body"`
}

type errBody struct {
	Body string `json:"body"`
}

type retBody struct {
	Body string `json:"cleaned_body"`
}

func (cfg *ApiConfig) ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {
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
	w.WriteHeader(200)
	respBody := retBody{Body: cleanedResult}
	dat, _ := json.Marshal(respBody)
	w.Write([]byte(dat))
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
