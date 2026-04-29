package port

import (
	"context"

	"golang-train/internal/payment/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p domain.Payment) (domain.Payment, error)
}
