package ws

import (
	"net/http"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/KirillLich/toomuCh/internal/tmerrors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WSHandler struct {
	hub    *Hub
	logger *zap.Logger
	cfg    config.WSConfig
}

func NewHandler(hub *Hub, logger *zap.Logger, cfg config.WSConfig) *WSHandler {
	return &WSHandler{hub: hub, logger: logger, cfg: cfg}
}

// TODO: move all params to config
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	const op = "ws.WSHandler.ServeWS"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket err", zap.String("op", op), zap.Error(err))
		code, errStr := tmerrors.MapErrorHttpStatus(err)
		http.Error(w, errStr, code)
		return
	}

	client := NewClient(h.logger, h.cfg, conn, h.hub)
	h.hub.register <- client
	go client.readPump()
	go client.writePump()
}
