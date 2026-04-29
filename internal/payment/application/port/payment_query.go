package port

import (
	"context"

	"golang-train/internal/payment/application/dto"
)

type PaymentQuery interface {
	ListWithUserName(ctx context.Context, limit int) ([]dto.PaymentListItemDTO, error)
}
