package messages

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	UserAuthError          = "UserAuthError"
	ValidationErrorMsg     = "ValidationError"
	ChangePasswordError    = "ChangePasswordError"
	ChangePasswordErrorMsg = "ChangePasswordErrorMsg"
	ChangePasswordOK       = "ChangePasswordOK"
	ChangePasswordOKMsg    = "ChangePasswordOKMsg"

	ValidationErrRequired         = "validation_required"
	ValidationErrAlphanumeric     = "validation_is_alphanumeric"
	ValidationErrLengthOutOfRange = "validation_length_out_of_range"

	ValidationInternalErrNotString               = "ValidationInternalErrNotString"
	ValidationInternalErrNotEnoughDigits         = "ValidationInternalErrNotEnoughDigits"
	ValidationInternalErrNotEnoughCapitalLetters = "ValidationInternalErrNotEnoughCapitalLetters"
)

var Messages = map[string]i18n.Message{
	UserAuthError: {
		ID:          UserAuthError,
		Description: "Ошибка проверка подлинности пользователя с помощью текущего пароля",
		Other:       "Authentication error",
	},
	ValidationErrorMsg: {
		ID:          ValidationErrorMsg,
		Description: "Ошибка валидации запроса",
		Other:       "Validation failed",
	},
	ChangePasswordError: {
		ID:          ChangePasswordError,
		Description: "Ошибка изменения пароля",
		Other:       "Change password error",
	},
	ChangePasswordErrorMsg: {
		ID:          ChangePasswordErrorMsg,
		Description: "Не удалось изменить пароль",
		Other:       "Failed to change password ",
	},
	ChangePasswordOK: {
		ID:          ChangePasswordOK,
		Description: "Смена пароля завершена",
		Other:       "Password change completed",
	},
	ChangePasswordOKMsg: {
		ID:          ChangePasswordOKMsg,
		Description: "Смена пароля выполнена успешно",
		Other:       "Password change completed successfully",
	},

	ValidationErrRequired: {
		ID:          ValidationErrRequired,
		Description: "Валидация. Значение должно быть указано",
		Other:       "The {{.Value}} cannot be blank",
	},
	ValidationErrAlphanumeric: {
		ID:          ValidationErrAlphanumeric,
		Description: "Валидация. Значение должно содержать цифры и буквы латинского алфавита",
		Other:       "The {{.Value}} must contain English letters and digits only",
	},
	ValidationErrLengthOutOfRange: {
		ID:          ValidationErrLengthOutOfRange,
		Description: "Валидация. Длина значения должна быть между {{.min}} and {{.max}}",
		Other:       "The length must be between {{.min}} and {{.max}}",
	},

	ValidationInternalErrNotString: {
		ID:          ValidationInternalErrNotString,
		Description: "Валидация. Значение должно быть строкой",
		Other:       "The {{.Value}} must be a string",
	},
	ValidationInternalErrNotEnoughDigits: {
		ID:          ValidationInternalErrNotEnoughDigits,
		Description: "Валидация. Значение должно содержать не менее {{.Digit}} цифры",
		One:         "The {{.Value}} must contain at least {{.Digit}} digit",
		Two:         "The {{.Value}} must contain at least {{.Digit}} digits",
		Few:         "The {{.Value}} must contain at least {{.Digit}} digits",
		Many:        "The {{.Value}} must contain at least {{.Digit}} digits",
		Other:       "The {{.Value}} must contain at least {{.Digit}} digits",
	},
	ValidationInternalErrNotEnoughCapitalLetters: {
		ID:          ValidationInternalErrNotEnoughCapitalLetters,
		Description: "Валидация. Значение должно содержать не менее {{.Digit}} заглавной буквы",
		One:         "The {{.Value}} must contain at least {{.Digit}} capital letter",
		Two:         "The {{.Value}} must contain at least {{.Digit}} capital letters",
		Few:         "The {{.Value}} must contain at least {{.Digit}} capital letters",
		Many:        "The {{.Value}} must contain at least {{.Digit}} capital letters",
		Other:       "The {{.Value}} must contain at least {{.Digit}} capital letters",
	},
}
