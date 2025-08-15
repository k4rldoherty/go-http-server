package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/k4rldoherty/go-http-server/internal/auth"
)

type polkaEventData struct {
	UserID string `json:"user_id"`
}

type polkaEventReq struct {
	Event string         `json:"event"`
	Data  polkaEventData `json:"data"`
}

func (cfg *ServerConfig) PolkaSubscriptionHandler(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("POLKA - %v", err)
		w.WriteHeader(401)
		return
	}
	if apiKey != cfg.PolkaKey {
		w.WriteHeader(401)
		return
	}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("POLKA - %v", err)
		w.WriteHeader(204)
		return
	}
	e := polkaEventReq{}
	err = json.Unmarshal(reqBody, &e)
	if err != nil {
		log.Printf("POLKA - %v", err)
		w.WriteHeader(204)
		return
	}

	if e.Event != "user.upgraded" {
		log.Println("POLKA - Event is not user.upgraded")
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(e.Data.UserID)
	if err != nil {
		log.Printf("POLKA - %v", err)
		w.WriteHeader(401)
		return
	}

	err = cfg.DBQueries.SubscribeUser(req.Context(), userID)
	if err != nil {
		log.Printf("POLKA - %v", err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}
