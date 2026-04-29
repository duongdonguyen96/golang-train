package usecase

import (
	"context"

	"golang-train/internal/payment/application/dto"
	"golang-train/internal/payment/application/port"
)

type ListPaymentsUsecase struct {
	q port.PaymentQuery
}

func NewListPaymentsUsecase(q port.PaymentQuery) *ListPaymentsUsecase {
	return &ListPaymentsUsecase{q: q}
}

func (uc *ListPaymentsUsecase) Execute(ctx context.Context, limit int) ([]dto.PaymentListItemDTO, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return uc.q.ListWithUserName(ctx, limit)
}
