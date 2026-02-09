package service

import (
	"context"
	"fmt"
	"time"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/KirillLich/toomuCh/internal/model"
	"github.com/KirillLich/toomuCh/internal/repository"
	"github.com/KirillLich/toomuCh/internal/tmerrors"
	"go.uber.org/zap"
)

type MessageCursor struct {
	PrevTime time.Time
	PrevId   int
}

type MessageService interface {
	CreateMessage(ctx context.Context, title string, text string) (*model.Message, error)
	GetMessage(ctx context.Context, limit int, cursor *MessageCursor) ([]model.Message, error)
	GetLatest(ctx context.Context, limit int) ([]model.Message, error)
	GetBefore(ctx context.Context, prevTime time.Time, prevId, limit int) ([]model.Message, error)
}

type messageService struct {
	logger  *zap.Logger
	repo    repository.MessageRepository
	app     config.AppConfig
	cleaner Cleaner
}

func NewMessageService(logger *zap.Logger, repo repository.MessageRepository, cleaner Cleaner, app config.AppConfig) *messageService {
	return &messageService{logger: logger, repo: repo, cleaner: cleaner, app: app}
}

func (s *messageService) CreateMessage(
	ctx context.Context,
	title string,
	text string,
) (*model.Message, error) {
	const op = "MessageService.CreateMessage"
	if len(title) > s.app.MessageMaxLen {
		return nil, fmt.Errorf("%s: %w", op, tmerrors.ErrTooLongTitle)
	}
	if len(text) > s.app.MessageMaxLen {
		return nil, fmt.Errorf("%s: %w", op, tmerrors.ErrTooLongText)
	}

	m := &model.Message{
		Title:     title,
		Text:      text,
		CreatedAt: time.Now().UTC(),
	}

	id, err := s.repo.CreateMessage(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	m.ID = id
	return m, nil
}

func (s *messageService) GetMessage(ctx context.Context, limit int, cursor *MessageCursor) ([]model.Message, error) {
	if cursor == nil {
		return s.GetLatest(ctx, limit)
	}

	return s.GetBefore(ctx, cursor.PrevTime, cursor.PrevId, limit)
}

func (s *messageService) GetLatest(ctx context.Context, limit int) ([]model.Message, error) {
	const op = "MessageService.GetLatest"
	if limit <= 0 {
		return nil, fmt.Errorf("%s: %w", op, tmerrors.ErrInvalidLimit)
	}

	return s.repo.GetLatest(ctx, limit)
}

func (s *messageService) GetBefore(ctx context.Context, prevTime time.Time, prevId, limit int) ([]model.Message, error) {
	const op = "MessageService.GetBefore"
	if limit <= 0 {
		return nil, fmt.Errorf("%s: %w", op, tmerrors.ErrInvalidLimit)
	}
	if prevId <= 0 {
		return nil, fmt.Errorf("%s: %w", op, tmerrors.ErrInvalidId)
	}

	return s.repo.GetBefore(ctx, prevTime, prevId, limit)
}
