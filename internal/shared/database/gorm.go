package database

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ctxKey struct{}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey{}, db)
}

func FromContext(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(ctxKey{}).(*gorm.DB)
	if !ok || db == nil {
		return nil, fmt.Errorf("db not found in context")
	}
	return db, nil
}

func openMySQL(host, port, user, pass, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&multiStatements=true", user, pass, host, port, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
