package domain

import "context"

type CheckPasswordI interface {
	Handle(ctx context.Context, login, password string) error
}

type ChangePasswordI interface {
	Handle(ctx context.Context, login, password string) error
}
