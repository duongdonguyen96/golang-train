package usecase

import (
	"context"

	"golang-train/internal/identity/application/port"
	"golang-train/internal/identity/domain"
)

type GetUserUsecase struct {
	repo port.UserRepository
}

func NewGetUserUsecase(repo port.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{repo: repo}
}

func (uc *GetUserUsecase) Execute(ctx context.Context, id uint64) (domain.User, error) {
	return uc.repo.GetByID(ctx, id)
}
