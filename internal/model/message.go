package model

import "time"

type Message struct {
	ID        int
	Text      string
	Title     string
	CreatedAt time.Time
}
