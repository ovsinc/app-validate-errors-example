package ports

import (
	"context"
	"errors"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/ovsinc/app-validate-errors-example/internal/service/domain"
)

const (
	_old_pass = "old_pass"
	_new_pass = "new_pass"
	_login    = "login"

	UserAuthError       = "UserAuthError"
	ValidationError     = "ValidationError"
	ChangePasswordError = "ChangePasswordError"
	ChangePasswordOK    = "ChangePasswordOK"
)

var (
	ErrAuthFailed           = errors.New("old pass auth failed")
	ErrChangePasswordFailed = errors.New("change password failed")
	ErrValidationError      = errors.New("request validation error")

	MsgUserAuthError = i18n.Message{
		ID:          UserAuthError,
		Description: "Ошибка проверка подлинности пользователя с помощью текущего пароля",
		Other:       "Authentication error",
	}
	MsgValidationErrorMsg = i18n.Message{
		ID:          ValidationError,
		Description: "Ошибка валидации запроса",
		Other:       "Validation failed",
	}
	MsgChangePasswordError = i18n.Message{
		ID:          ChangePasswordError,
		Description: "Не удалось изменить пароль",
		Other:       "Failed to change password",
	}
	MsgChangePasswordOK = i18n.Message{
		ID:          ChangePasswordOK,
		Description: "Смена пароля выполнена успешно",
		Other:       "Password change completed successfully",
	}
)

type ChangePasswordRequest struct {
	Login       string `json:"login"`
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	Lang        string `json:"lang"`
}

type ChangePasswordResponse struct {
	Success bool         `json:"success"`
	Payload Payload      `json:"payload"`
	Error   ErrorPayload `json:"errors"`
}

type httpServer struct {
	appChanage domain.ChangePasswordI
	appCheck   domain.CheckPasswordI
	bundle     *i18n.Bundle
}

func NewHttpServer(
	appChanage domain.ChangePasswordI,
	appCheck domain.CheckPasswordI,
	bundle *i18n.Bundle,
) PasswordChange {
	return &httpServer{
		appChanage: appChanage,
		appCheck:   appCheck,
		bundle:     bundle,
	}
}

func (h *httpServer) ChangePassword(
	ctx context.Context,
	req *ChangePasswordRequest,
) (*ChangePasswordResponse, int, error) {

	localizer := i18n.NewLocalizer(h.bundle, req.Lang)

	if err := h.appCheck.Handle(ctx, req.Login, req.OldPassword); err != nil {
		msg := "Authentication error"
		if lmsg, lerr := localizer.LocalizeMessage(&MsgUserAuthError); lerr == nil {
			msg = lmsg
		}
		return &ChangePasswordResponse{
				Success: false,
				Error: ErrorPayload{
					_old_pass: []string{msg},
				},
			},
			http.StatusBadRequest,
			ErrAuthFailed
	}

	if err := req.Validate(); err != nil {
		return &ChangePasswordResponse{
				Success: false,
				Error:   ValidatorErrors(err, localizer),
			},
			http.StatusBadRequest,
			ErrValidationError
	}

	if err := h.appChanage.Handle(ctx, req.Login, req.Password); err != nil {
		merr := "Change password error"
		if lmsg, lerr := localizer.LocalizeMessage(&MsgChangePasswordError); lerr == nil {
			merr = lmsg
		}
		return &ChangePasswordResponse{
				Success: false,
				Error: ErrorPayload{
					_new_pass: []string{merr},
				},
			},
			http.StatusInternalServerError,
			ErrChangePasswordFailed
	}

	msg := "Password change completed successfully"
	if lmsg, lerr := localizer.LocalizeMessage(&MsgChangePasswordOK); lerr == nil {
		msg = lmsg
	}
	return &ChangePasswordResponse{
			Success: true,
			Payload: Payload{
				Status: msg,
			},
		},
		http.StatusOK,
		nil
}
