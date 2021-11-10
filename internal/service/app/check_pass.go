package app

import (
	"context"

	"github.com/ovsinc/app-validate-errors-example/internal/service/domain"
)

func NewCheckPassword(repo domain.Repository) domain.CheckPasswordI {
	return &checkPassword{repo: repo}
}

type checkPassword struct {
	repo domain.Repository
}

func (chp *checkPassword) Handle(ctx context.Context, login, password string) error {
	return chp.repo.ChangePassword(ctx, login, password)
}
