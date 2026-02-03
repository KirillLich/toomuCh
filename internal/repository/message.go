package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/KirillLich/toomuCh/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, m *model.Message) (int, error)
	GetLatest(ctx context.Context, limit int) ([]model.Message, error)
	GetBefore(ctx context.Context, prevTime time.Time, prevId, limit int) ([]model.Message, error)
	DeleteBefore(ctx context.Context, delTime time.Time) (int, error)
}

type messageRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewMessageRepo(logger *zap.Logger, db *gorm.DB) *messageRepository {
	return &messageRepository{logger: logger, db: db}
}

func (repo *messageRepository) CreateMessage(ctx context.Context, m *model.Message) (int, error) {
	const op = "MessageRepository.CreateMessage"
	result := repo.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}
	return m.ID, nil
}

func (repo *messageRepository) GetLatest(ctx context.Context, limit int) ([]model.Message, error) {
	const op = "MessageRepository.GetLatest"
	var messages []model.Message
	result := repo.db.WithContext(ctx).Order("created_at DESC, id DESC").Limit(limit).Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}
	return messages, nil
}

func (repo *messageRepository) GetBefore(ctx context.Context, prevTime time.Time, prevId, limit int) ([]model.Message, error) {
	const op = "MessageRepository.GetBefore"
	var messages []model.Message
	tempGormDB := repo.db.WithContext(ctx).Order("created_at DESC, id DESC").Where("created_at < ? OR (created_at = ? AND id < ?)", prevTime, prevTime, prevId)
	result := tempGormDB.Limit(limit).Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}
	return messages, nil
}

func (repo *messageRepository) DeleteBefore(ctx context.Context, delTime time.Time) (int, error) {
	const op = "MessageRepository.DeleteBefore"
	tempGormDB := repo.db.WithContext(ctx).Where("created_at <= ?", delTime)
	result := tempGormDB.Delete(&model.Message{})
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}
	return int(result.RowsAffected), nil
}
