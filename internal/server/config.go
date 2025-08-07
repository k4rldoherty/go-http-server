// Package server
package server

import (
	"sync/atomic"

	"github.com/k4rldoherty/go-http-server/internal/auth"
	"github.com/k4rldoherty/go-http-server/internal/database"
)

type ServerConfig struct {
	FileServerHits atomic.Int32
	DBQueries      *database.Queries
	Platform       string
	JWTCfg         *auth.JWTConfig
}
