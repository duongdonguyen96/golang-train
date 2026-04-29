package domain

import "time"

type Payment struct {
	ID        uint64
	UserID    uint64
	Amount    int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
