package adapter

import (
	"context"

	"github.com/ovsinc/app-validate-errors-example/internal/service/domain"
)

/*
type Repository interface {
	ChangePassword(ctx context.Context, login, password string) error
	CheckPassord(ctx context.Context, login, password string) error
}
*/

type passwordRepository struct{}

func NewPasswordRepository() domain.Repository {
	return &passwordRepository{}
}

func (pr *passwordRepository) ChangePassword(ctx context.Context, login, password string) error {
	return nil
}

func (pr *passwordRepository) CheckPassord(ctx context.Context, login, password string) error {
	return nil
}
