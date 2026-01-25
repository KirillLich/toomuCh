package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KirillLich/toomuCh/internal/config"
	"github.com/KirillLich/toomuCh/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

func InitPostgres(cfg config.DBConfig, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode)

	gormLogger := gl.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gl.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gl.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return &gorm.DB{}, fmt.Errorf("error connerction db: %w", err)
	}

	if err = db.AutoMigrate(&model.Message{}); err != nil {
		return &gorm.DB{}, fmt.Errorf("error while migration: %w", err)
	}

	return db, nil
}
