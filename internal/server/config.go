// Package server
package server

import (
	"sync/atomic"

	"github.com/k4rldoherty/go-http-server/internal/database"
)

type APIConfig struct {
	FileServerHits atomic.Int32
	DBQueries      *database.Queries
}
