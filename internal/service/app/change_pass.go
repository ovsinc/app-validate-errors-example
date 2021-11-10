package app

import (
	"context"

	"github.com/ovsinc/app-validate-errors-example/internal/service/domain"
)

func NewChangePassword(repo domain.Repository) domain.ChangePasswordI {
	return &changePassword{repo: repo}
}

type changePassword struct {
	repo domain.Repository
}

func (chp *changePassword) Handle(ctx context.Context, login, password string) error {
	return chp.repo.ChangePassword(ctx, login, password)
}
