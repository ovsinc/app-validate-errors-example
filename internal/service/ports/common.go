package ports

import (
	"context"
)

type Payload struct {
	Status string `json:"status"`
}

type ErrorOut struct {
	Message   string                 `json:"message"`
	Operation string                 `json:"operation"`
	ID        string                 `json:"id"`
	Context   map[string]interface{} `json:"context"`
}

type ErrorPayload map[string][]string

//

type PasswordChange interface {
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, int, error)
}
