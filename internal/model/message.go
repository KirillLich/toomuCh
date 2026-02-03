package model

import "time"

type Message struct {
	ID        int `gorm:"primaryKey"`
	Text      string
	Title     string
	CreatedAt time.Time `gorm:"index"`
}
