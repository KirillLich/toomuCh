package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/KirillLich/toomuCh/internal/db"
	"github.com/KirillLich/toomuCh/internal/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	cfg := config.SetConfig()
	logger := logger.NewLogger(cfg.App.LogLVL, cfg.Env)
	logger.Info("Config fields",
		zap.String("env", cfg.Env),

		zap.String("server.host", cfg.Server.Host),
		zap.Int("server.port", cfg.Server.Port),

		zap.String("db.host", cfg.DB.Host),
		zap.Int("db.port", cfg.DB.Port),
		zap.String("db.name", cfg.DB.Name),

		zap.Duration("app.ttl", cfg.App.TTL),
		zap.Int("app.messageMaxLen", cfg.App.MessageMaxLen),
	)

	_, err := db.InitPostgres(cfg.DB, logger)
	if err != nil {
		logger.Fatal("error initialisation db", zap.Error(err))
	}

	r := chi.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			HealthResponse{Status: fmt.Sprintf(
				"Hello from main. Listen and serve at: %s:%d",
				cfg.Server.Host,
				cfg.Server.Port,
			)})
	})

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r)
}
