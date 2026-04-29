package persistence

import (
	"context"

	"golang-train/internal/payment/domain"
	"golang-train/internal/shared/database"
)

type paymentModel struct {
	ID       uint64 `gorm:"primaryKey"`
	UserID   uint64 `gorm:"index"`
	Amount   int64
	Currency string `gorm:"size:8"`
}

func (paymentModel) TableName() string { return "payments" }

type GormPaymentRepository struct{}

func NewGormPaymentRepository() *GormPaymentRepository { return &GormPaymentRepository{} }

func (r *GormPaymentRepository) Create(ctx context.Context, p domain.Payment) (domain.Payment, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return domain.Payment{}, err
	}
	if err := db.AutoMigrate(&paymentModel{}); err != nil {
		return domain.Payment{}, err
	}
	m := paymentModel{UserID: p.UserID, Amount: p.Amount, Currency: p.Currency}
	if err := db.Create(&m).Error; err != nil {
		return domain.Payment{}, err
	}
	p.ID = m.ID
	return p, nil
}
