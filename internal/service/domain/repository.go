package domain

import (
	"context"
)

type Repository interface {
	ChangePassword(ctx context.Context, login, password string) error
	CheckPassord(ctx context.Context, login, password string) error
}
