package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"golang-train/internal/auth/application/port"
)

type LoginUsecase struct {
	q         port.UserCredentialsQuery
	jwtSecret []byte
}

func NewLoginUsecase(q port.UserCredentialsQuery, jwtSecret []byte) *LoginUsecase {
	return &LoginUsecase{q: q, jwtSecret: jwtSecret}
}

func (uc *LoginUsecase) Execute(ctx context.Context, tenant string, email string, password string) (string, error) {
	// tenant is validated by the tenant middleware header resolver.
	userID, hash, err := uc.q.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials..... khong co data cho email nay: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":    fmt.Sprintf("%d", userID),
		"tenant": tenant,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(uc.jwtSecret)
}
