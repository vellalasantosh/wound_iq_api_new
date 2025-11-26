package handlers

import (
	"database/sql"

	"github.com/vellalasantosh/wound_iq_api_new/internal/config"
	"go.uber.org/zap"
)

type Handlers struct {
	DB  *sql.DB
	Log *zap.Logger
	Cfg *config.Config
}

func NewHandlers(db *sql.DB, log *zap.Logger, cfg *config.Config) *Handlers {
	return &Handlers{
		DB:  db,
		Log: log,
		Cfg: cfg,
	}
}
