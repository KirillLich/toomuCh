package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KirillLich/toomuCh/internal/service"
	"github.com/KirillLich/toomuCh/internal/tmerrors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type MessageHandler interface {
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetLatest(w http.ResponseWriter, r *http.Request)
}

type messageHandler struct {
	serv   service.MessageService
	logger *zap.Logger
}

type createMessageRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func NewMessageHandler(serv service.MessageService, logger *zap.Logger) *messageHandler {
	return &messageHandler{serv: serv, logger: logger}
}

func (h *messageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	const op = "MessageHandler.CreateMessage"
	var req createMessageRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		code, errStr := http.StatusBadRequest, "invalid request body"
		http.Error(w, errStr, code)
		return
	}

	msg, err := h.serv.CreateMessage(r.Context(), req.Title, req.Text)
	if err != nil {
		code, errStr := tmerrors.MapErrorHttpStatus(err)
		http.Error(w, errStr, code)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, msg)
}

func (h *messageHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	const op = "MessageHandler.GetLatest"
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidLimit)
		http.Error(w, errStr, code)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidLimit)
		http.Error(w, errStr, code)
		return
	}

	timeStr := r.URL.Query().Get("cursor_time")
	idStr := r.URL.Query().Get("cursor_id")

	var cursor *service.MessageCursor

	if timeStr == "" && idStr == "" {
		cursor = nil
	} else if timeStr == "" || idStr == "" {
		if timeStr == "" {
			code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidTime)
			http.Error(w, errStr, code)
			return
		}
		if idStr == "" {
			code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidId)
			http.Error(w, errStr, code)
			return
		}
	} else {
		time, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidTime)
			http.Error(w, errStr, code)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			code, errStr := tmerrors.MapErrorHttpStatus(tmerrors.ErrInvalidId)
			http.Error(w, errStr, code)
			return
		}

		cursor.PrevId = id
		cursor.PrevTime = time
	}

	msgs, err := h.serv.GetMessage(r.Context(), limit, cursor)
	if err != nil {
		code, errStr := tmerrors.MapErrorHttpStatus(err)
		http.Error(w, errStr, code)
		return
	}

	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, msgs)
}
