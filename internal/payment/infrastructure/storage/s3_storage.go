package storage

import (
	"context"
	"fmt"

	"golang-train/internal/shared/config"
)

// S3LikeStorage is a minimal adapter returning a URL-like string.
// It is intentionally simple to keep the example compiling without AWS SDK.
// Replace with AWS SDK / MinIO client in real systems.

type S3LikeStorage struct {
	cfg config.S3Config
}

func NewS3LikeStorage(cfg config.S3Config) *S3LikeStorage {
	return &S3LikeStorage{cfg: cfg}
}

func (s *S3LikeStorage) PutObject(ctx context.Context, key string, contentType string, data []byte) (string, error) {
	_ = ctx
	_ = contentType
	_ = data
	// For demo: we don't actually upload. Return a deterministic URL.
	return fmt.Sprintf("s3://%s/%s", s.cfg.Bucket, key), nil
}
