package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/k4rldoherty/go-http-server/internal/auth"
)

type refreshRes struct {
	Token string `json:"token"`
}

func (cfg *ServerConfig) RefreshHandler(w http.ResponseWriter, req *http.Request) {
	t, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("REFRESH GET BEARER - %v", err)
		w.WriteHeader(401)
		return
	}
	userID, err := cfg.DBQueries.GetUserFromRefreshToken(req.Context(), t)
	if err != nil {
		log.Printf("REFRESH GET USER FROM REFRESH - %v", err)
		w.WriteHeader(401)
		return
	}
	// Generate a new access token
	newAccessToken, err := auth.MakeJWT(userID.UUID, cfg.JWTCfg)
	if err != nil {
		log.Printf("REFRESH MAKE JWT - %v", err)
		w.WriteHeader(401)
		return
	}
	// otherwise return 200 with res
	res := refreshRes{
		Token: newAccessToken,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("REFRESH MARSHAL JSON - %v", err)
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(200)
	w.Write(resJSON)
}
