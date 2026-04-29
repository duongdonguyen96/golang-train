package usecase

import (
	"context"
	"fmt"

	"golang-train/internal/payment/application/port"
	"golang-train/internal/payment/domain"
)

type CreatePaymentUsecase struct {
	repo port.PaymentRepository
}

func NewCreatePaymentUsecase(repo port.PaymentRepository) *CreatePaymentUsecase {
	return &CreatePaymentUsecase{repo: repo}
}

func (uc *CreatePaymentUsecase) Execute(ctx context.Context, userID uint64, amount int64, currency string) (domain.Payment, error) {
	if userID == 0 {
		return domain.Payment{}, fmt.Errorf("user_id required")
	}
	if amount <= 0 {
		return domain.Payment{}, fmt.Errorf("amount must be positive")
	}
	if currency == "" {
		currency = "USD"
	}
	return uc.repo.Create(ctx, domain.Payment{UserID: userID, Amount: amount, Currency: currency})
}
