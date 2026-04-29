package persistence

import (
	"context"
	"errors"
	"fmt"

	"golang-train/internal/identity/domain"
	"golang-train/internal/shared/database"

	"gorm.io/gorm"
)

type userModel struct {
	ID           uint64 `gorm:"primaryKey"`
	Email        string `gorm:"size:255;uniqueIndex"`
	PasswordHash string `gorm:"size:255"`
	Name         string `gorm:"size:255"`
	CreatedAt    int64
	UpdatedAt    int64
}

func (userModel) TableName() string { return "users" }

type GormUserRepository struct{}

func NewGormUserRepository() *GormUserRepository { return &GormUserRepository{} }

func (r *GormUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	m := userModel{Email: u.Email, PasswordHash: u.PasswordHash, Name: u.Name}
	if err := db.AutoMigrate(&userModel{}); err != nil {
		return domain.User{}, err
	}
	if err := db.Create(&m).Error; err != nil {
		return domain.User{}, err
	}
	u.ID = m.ID
	return u, nil
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	if err := db.AutoMigrate(&userModel{}); err != nil {
		return domain.User{}, err
	}

	var m userModel
	if err := db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, err
	}

	return domain.User{ID: m.ID, Email: m.Email, PasswordHash: m.PasswordHash, Name: m.Name}, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	if err := db.AutoMigrate(&userModel{}); err != nil {
		return domain.User{}, err
	}

	var m userModel
	if err := db.First(&m, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, err
	}

	return domain.User{ID: m.ID, Email: m.Email, PasswordHash: m.PasswordHash, Name: m.Name}, nil
}

// ---- Auth query adapter (port lives in auth module)

type userCredentialDTO struct {
	ID           uint64
	PasswordHash string
}

type GormUserCredentialQuery struct{}

func NewGormUserCredentialQuery() *GormUserCredentialQuery { return &GormUserCredentialQuery{} }

func (q *GormUserCredentialQuery) GetCredentialsByEmail(ctx context.Context, email string) (uint64, string, error) {
	db, err := database.FromContext(ctx)
	if err != nil {
		return 0, "", err
	}
	// if err := db.AutoMigrate(&userModel{}); err != nil {
	// 	return 0, "", err
	// }

	var out userCredentialDTO
	// Query builder avoids returning ORM model outside infrastructure.
	if err := db.Table("users").Select("id, password_hash").Where("email = ?", email).Scan(&out).Error; err != nil {
		return 0, "", err
	}
	if out.ID == 0 {
		return 0, "", fmt.Errorf("invalid credentials")
	}
	return out.ID, out.PasswordHash, nil
}
