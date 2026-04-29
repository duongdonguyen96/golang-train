package query

import (
	"context"

	"golang-train/internal/payment/application/dto"
	"golang-train/internal/shared/database"
)

type GormPaymentQuery struct{}

func NewGormPaymentQuery() *GormPaymentQuery { return &GormPaymentQuery{} }

func (q *GormPaymentQuery) ListWithUserName(ctx context.Context, limit int) ([]dto.PaymentListItemDTO, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var out []dto.PaymentListItemDTO

	// Read model: JOIN via query builder, no Preload, no ORM model exposure.
	err = db.Table("payments").
		Select("payments.id as payment_id, payments.user_id, users.name as user_name, payments.amount, payments.currency, payments.created_at").
		Joins("JOIN users ON users.id = payments.user_id").
		Order("payments.id desc").
		Limit(limit).
		Scan(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}
