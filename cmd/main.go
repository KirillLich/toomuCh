package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/KirillLich/toomuCh/internal/db"
	"github.com/KirillLich/toomuCh/internal/handler"
	"github.com/KirillLich/toomuCh/internal/logger"
	"github.com/KirillLich/toomuCh/internal/repository"
	"github.com/KirillLich/toomuCh/internal/service"
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

	conn, err := db.InitPostgres(cfg.DB, logger)
	if err != nil {
		logger.Fatal("error initialisation db", zap.Error(err))
	}

	repo := repository.NewMessageRepo(logger, conn)
	cl := service.NewCleaner(repo, cfg.App.TTL, cfg.App.SleepTime, logger, context.Background())
	serv := service.NewMessageService(logger, repo, cl, cfg.App)
	handler := handler.NewMessageHandler(serv, logger)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			HealthResponse{Status: fmt.Sprintf(
				"Hello from main. Listen and serve at: %s:%d",
				cfg.Server.Host,
				cfg.Server.Port,
			)})
	})
	r.Post("/message", handler.CreateMessage)
	r.Get("/message", handler.GetLatest)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r)
}
