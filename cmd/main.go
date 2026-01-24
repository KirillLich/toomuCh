package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/go-chi/chi"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	cfg := config.SetConfig()
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

	fmt.Println(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	fmt.Println(cfg)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r)
}
