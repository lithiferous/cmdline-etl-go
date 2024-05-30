package api

import (
	db "github.com/lithiferous/cmd-etl/db/sqlc"
	"github.com/lithiferous/cmd-etl/util"
)

type AppState struct{}

// State serves as an app state for our snapshot service.
type State struct {
	Config util.Config
	Store  db.Store
}

// NewServer creates a new configured server.
// Later on will validate config here.
func NewState(config util.Config, store db.Store) (*State, error) {

	server := &State{
		Config: config,
		Store:  store,
	}

	return server, nil
}
