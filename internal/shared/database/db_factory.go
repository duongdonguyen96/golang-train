package database

import (
	"context"
	"fmt"
	"sync"

	"golang-train/internal/shared/config"

	"gorm.io/gorm"
)

type Factory interface {
	ForTenant(ctx context.Context, tenant string) (*gorm.DB, error)
}

type mySQLFactory struct {
	cfg   config.DBConfig
	mu    sync.RWMutex
	cache map[string]*gorm.DB
}

func NewMySQLFactory(cfg config.DBConfig) Factory {
	return &mySQLFactory{cfg: cfg, cache: map[string]*gorm.DB{}}
}

func (f *mySQLFactory) ForTenant(ctx context.Context, tenant string) (*gorm.DB, error) {
	schema := fmt.Sprintf("tenant_%s", tenant)

	f.mu.RLock()
	if db := f.cache[schema]; db != nil {
		f.mu.RUnlock()
		return db, nil
	}
	f.mu.RUnlock()

	// Ensure schema exists using a bootstrap connection.
	bootstrap, err := openMySQL(f.cfg.Host, f.cfg.Port, f.cfg.User, f.cfg.Password, "mysql")
	if err != nil {
		return nil, err
	}
	if err := bootstrap.Exec("CREATE DATABASE IF NOT EXISTS `" + schema + "`").Error; err != nil {
		return nil, err
	}

	db, err := openMySQL(f.cfg.Host, f.cfg.Port, f.cfg.User, f.cfg.Password, schema)
	if err != nil {
		return nil, err
	}

	f.mu.Lock()
	f.cache[schema] = db
	f.mu.Unlock()

	return db, nil
}
