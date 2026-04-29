package port

import "context"

type UserCredentialsQuery interface {
	GetCredentialsByEmail(ctx context.Context, email string) (userID uint64, passwordHash string, err error)
}
