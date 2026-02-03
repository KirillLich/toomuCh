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

// type MessageRepository interface {
// 	CreateMessage(ctx context.Context, m *model.Message) (int, error)
// 	GetLatest(ctx context.Context, limit int) ([]model.Message, error)
// 	GetBefore(ctx context.Context, prevTime time.Time, prevId, limit int) ([]model.Message, error)
// 	DeleteBefore(ctx context.Context, delTime time.Time) (int, error)
// }

type MessageService interface {
	CreateMessage(ctx context.Context, m *model.Message) (int, error)
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

func (s *messageService) CreateMessage(ctx context.Context, m *model.Message) (int, error) {
	const op = "MessageService.CreateMessage"
	if len(m.Title) > s.app.MessageMaxLen {
		return 0, fmt.Errorf("%s: %w", op, tmerrors.ErrTooLongTitle)
	}
	if len(m.Text) > s.app.MessageMaxLen {
		return 0, fmt.Errorf("%s: %w", op, tmerrors.ErrTooLongText)
	}

	return s.repo.CreateMessage(ctx, m)
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
