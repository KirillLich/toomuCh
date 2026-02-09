package service

import (
	"context"
	"fmt"
	"time"

	"github.com/KirillLich/toomuCh/internal/repository"
	"go.uber.org/zap"
)

type Cleaner interface {
	run()
}

type cleaner struct {
	repo      repository.MessageRepository
	ttl       time.Duration
	sleepTime time.Duration
	logger    *zap.Logger
	ctx       context.Context
}

func (c *cleaner) run() {
	const op = "Cleaner.Run"
	ticker := time.NewTicker(c.sleepTime)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			cutOff := time.Now().Add(-c.ttl)
			num, err := c.repo.DeleteBefore(c.ctx, cutOff)
			if err != nil {
				c.logger.Error("TTL cleanup error", zap.Error(fmt.Errorf("%s: %w", op, err)))
			} else {
				c.logger.Info("TTL cleanup succesfully finished", zap.Int("rows deleted", num))
			}
		}
	}
}

func NewCleaner(repo repository.MessageRepository,
	ttl time.Duration,
	sleepTime time.Duration,
	logger *zap.Logger,
	ctx context.Context,
) *cleaner {
	c := cleaner{repo: repo, ttl: ttl, sleepTime: sleepTime, logger: logger, ctx: ctx}
	go c.run()
	return &c
}
