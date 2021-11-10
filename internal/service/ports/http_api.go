package ports

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/ovsinc/errors"

	"github.com/ovsinc/app-validate-errors-example/internal/service/domain"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ChangePasswordRequest struct {
	Login       string `json:"login"`
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	Lang        string `json:"lang"`
}

func (req *ChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.Login,
			is.Alphanumeric, validation.Required, validation.Length(2, 0)),
		validation.Field(&req.Password,
			validation.Required, validation.Length(5, 0), validation.By(checkSimplePass)),
	)
}

var (
	ErrNotString        = errors.New("value is not a string")
	ErrNoDigits         = errors.New("value not contain digits")
	ErrNoCapitalLetters = errors.New("value not contain capital letters")
)

const (
	_minDigit   = 1
	_minCapital = 1
)

func checkSimplePass(value interface{}) error {
	var (
		s  string
		ok bool
	)
	if s, ok = value.(string); !ok {
		return ErrNotString
	}

	var digits, alpha, capital int

	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r, _ := utf8.DecodeRune(scanner.Bytes())
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r >= 'a' && r <= 'z':
			alpha++
		case r >= 'A' && r <= 'Z':
			capital++
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	switch {
	case digits < _minDigit:
		return ErrNoDigits
	case capital < _minCapital:
		return ErrNoCapitalLetters
	}

	return nil
}

type ChangePasswordResponse struct {
	Common
	Payload Payload      `json:"payload"`
	Error   ErrorPayload `json:"errors"`
}

type PasswordChange interface {
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, int, error)
}

type httpServer struct {
	appChanage domain.ChangePasswordI
	appCheck   domain.CheckPasswordI
}

func NewHttpServer(
	appChanage domain.ChangePasswordI,
	appCheck domain.CheckPasswordI) PasswordChange {
	return &httpServer{
		appChanage: appChanage,
		appCheck:   appCheck,
	}
}

func ValidatorErrors(err error) map[string][]string {
	const (
		_op     = "validation"
		_common = "common"
	)

	log.Printf("[INFO] validation error: %v", err)

	errFields := make(map[string][]string)

	es, ok := err.(validation.Errors)
	if !ok {
		log.Println("[WARN] validation: not a ozzo validation error")
		errFields[_common] = []string{err.Error()}
		return errFields
	}

	out := make([]string, 0, len(es))

	for i, err := range es {
		em, ok := err.(validation.Error)
		if !ok {
			out = append(out, err.Error())
		} else {
			out = append(
				out,
				fmt.Sprintf("%s: %s", em.Code(), em.Error()),
			)
		}
		errFields[i] = out
	}

	return errFields
}

var (
	ErrAuthFailed           = errors.New("old pass auth failed")
	ErrChangePasswordFailed = errors.New("change password failed")
	ErrValidationError      = errors.New("request validation error")
)

const (
	_old_pass = "old_pass"
	_new_pass = "new_pass"
	_login    = "login"
)

func (h *httpServer) ChangePassword(
	ctx context.Context,
	req *ChangePasswordRequest,
) (*ChangePasswordResponse, int, error) {
	if err := h.appCheck.Handle(ctx, req.Login, req.OldPassword); err != nil {
		return &ChangePasswordResponse{
				Common: Common{
					Success: false,
					Message: "Ошибка проверка подлинности пользователя",
				},
				Error: ErrorPayload{
					_old_pass: []string{"Ошибка аутентификации"},
				},
			},
			http.StatusBadRequest,
			ErrAuthFailed
	}

	if err := req.Validate(); err != nil {
		return &ChangePasswordResponse{
				Common: Common{
					Success: false,
					Message: "Ошибка валидации запроса",
				},
				Error: ValidatorErrors(err),
			},
			http.StatusBadRequest,
			ErrValidationError
	}

	if err := h.appChanage.Handle(ctx, req.Login, req.Password); err != nil {
		return &ChangePasswordResponse{
				Common: Common{
					Success: false,
					Message: "Ошибка изменения пароля",
				},
				Error: ErrorPayload{
					_new_pass: []string{"Не удалось изменить пароль"},
				},
			},
			http.StatusInternalServerError,
			ErrChangePasswordFailed
	}

	return &ChangePasswordResponse{
			Common: Common{
				Success: true,
				Message: "Смена пароля выполнена успешно",
			},
			Payload: Payload{
				Status: "Смена пароля выполнена успешно",
			},
		},
		http.StatusOK,
		nil
}
