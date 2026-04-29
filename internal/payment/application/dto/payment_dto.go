package dto

import "time"

type PaymentListItemDTO struct {
	PaymentID uint64    `json:"payment_id"`
	UserID    uint64    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
