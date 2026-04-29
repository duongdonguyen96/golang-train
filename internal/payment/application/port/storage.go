package port

import "context"

type Storage interface {
	PutObject(ctx context.Context, key string, contentType string, data []byte) (string, error)
}
