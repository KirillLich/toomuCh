package ws

import (
	"time"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const defaultTicker = 50 * time.Second

type Client struct {
	logger *zap.Logger
	cfg    config.WSConfig
	conn   *websocket.Conn
	hub    *Hub
	send   chan []byte
}

func NewClient(logger *zap.Logger,
	cfg config.WSConfig,
	conn *websocket.Conn,
	hub *Hub,
) *Client {
	client := &Client{
		logger: logger,
		cfg:    cfg,
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, cfg.BuffSize),
	}

	client.conn.SetPongHandler(func(appData string) error {
		return client.conn.SetReadDeadline(time.Now().Add(cfg.ReadDeadline))
	})

	return client
}

func (c *Client) readPump() {
	const op = "ws.Client.ReadPump"
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(c.cfg.ReadDeadline))

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("unexpected err", zap.String("op", op), zap.Error(err))
			}
			return
		}
	}
}

func (c *Client) writePump() {
	const op = "ws.Client.WritePump"

	var ticker *time.Ticker
	if c.cfg.PingPeriod != 0 {
		ticker = time.NewTicker(c.cfg.PingPeriod)
	} else {
		ticker = time.NewTicker(defaultTicker)
	}
	defer func() {
		c.hub.unregister <- c
		ticker.Stop()
		c.conn.Close()
	}()
	c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteDeadline))

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteDeadline))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.logger.Error("unexpected err", zap.String("op", op), zap.Error(err))
				}
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteDeadline))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.logger.Error("unexpected err", zap.String("op", op), zap.Error(err))
				}
				return
			}
		}

	}
}
