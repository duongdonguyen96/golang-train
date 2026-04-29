package port

import (
	"context"

	"golang-train/internal/identity/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByID(ctx context.Context, id uint64) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}
