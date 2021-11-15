package ports

import (
	"bufio"
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	_minDigit   = 1
	_minCapital = 1

	ValidationErrRequired         = "validation_required"
	ValidationErrAlphanumeric     = "validation_is_alphanumeric"
	ValidationErrLengthOutOfRange = "validation_length_out_of_range"

	ValidationInternalErrNotString               = "ValidationInternalErrNotString"
	ValidationInternalErrNotEnoughDigits         = "ValidationInternalErrNotEnoughDigits"
	ValidationInternalErrNotEnoughCapitalLetters = "ValidationInternalErrNotEnoughCapitalLetters"
)

var (
	ErrNotString               = errors.New("value is not a string")
	ErrNotEnoughDigits         = errors.New("there are not enough digits in the value")
	ErrNotEnoughCapitalLetters = errors.New("there are not enough capital letters in the value")

	MsgValidationErrRequired = i18n.Message{
		ID:          ValidationErrRequired,
		Description: "Валидация. Значение должно быть указано",
		Other:       "The value cannot be blank",
	}
	MsgValidationErrAlphanumeric = i18n.Message{
		ID:          ValidationErrAlphanumeric,
		Description: "Валидация. Значение должно содержать цифры и буквы латинского алфавита",
		Other:       "The value must contain English letters and digits only",
	}
	MsgValidationErrLengthOutOfRange = i18n.Message{
		ID:          ValidationErrLengthOutOfRange,
		Description: "Валидация. Длина значения должна быть между {{.min}} and {{.max}}",
		Other:       "The length must be between {{.min}} and {{.max}}",
	}

	MsgValidationInternalErrNotString = i18n.Message{
		ID:          ValidationInternalErrNotString,
		Description: "Валидация. Значение должно быть строкой",
		Other:       "The value must be a string",
	}
	MsgValidationInternalErrNotEnoughDigits = i18n.Message{
		ID:          ValidationInternalErrNotEnoughDigits,
		Description: "Валидация. Значение должно содержать не менее {{.digit}} цифры",
		One:         "The value must contain at least {{.digit}} digit",
		Two:         "The value must contain at least {{.digit}} digits",
		Few:         "The value must contain at least {{.digit}} digits",
		Many:        "The value must contain at least {{.digit}} digits",
		Other:       "The value must contain at least {{.digit}} digits",
	}
	MsgValidationInternalErrNotEnoughCapitalLetters = i18n.Message{
		ID:          ValidationInternalErrNotEnoughCapitalLetters,
		Description: "Валидация. Значение должно содержать не менее {{.digit}} заглавной буквы",
		One:         "The value must contain at least {{.digit}} capital letter",
		Two:         "The value must contain at least {{.digit}} capital letters",
		Few:         "The value must contain at least {{.digit}} capital letters",
		Many:        "The value must contain at least {{.digit}} capital letters",
		Other:       "The value must contain at least {{.digit}} capital letters",
	}
)

func ValidatorErrors(err error, localizer *i18n.Localizer) map[string][]string {
	const (
		_op     = "validation"
		_common = "common"
	)

	log.Printf("[INFO] validation error: %v", err)

	errFields := make(map[string][]string)

	if errs, ok := err.(validation.Errors); !ok {
		log.Println("[WARN] validation: not a ozzo validation error")
		errFields[_common] = []string{err.Error()}

	} else {
		for i, e := range errs {
			switch es := e.(type) {
			case validation.Error:
				msg := err.Error()
				if lmsg, lerr := localizer.Localize(&i18n.LocalizeConfig{
					MessageID:    es.Code(),
					TemplateData: es.Params(),
				}); lerr == nil {
					msg = lmsg
				}
				errFields[i] = []string{msg}
			case *checkPassError:
				msg := err.Error()
				if lmsg, lerr := localizer.Localize(&i18n.LocalizeConfig{
					MessageID: es.code,
				}); lerr == nil {
					msg = lmsg
				}
				errFields[i] = []string{msg}

			case *checkPassErrors:
				out := make([]string, 0, len(es.errs))
				for _, em := range es.errs {
					msg := err.Error()
					switch {
					case em.digits > 0:
						if lmsg, lerr := localizer.Localize(&i18n.LocalizeConfig{
							MessageID:   em.code,
							PluralCount: em.digits,
							TemplateData: map[string]int{
								"digit": em.digits,
							},
						}); lerr == nil {
							msg = lmsg
						}

					case em.capital > 0:
						if lmsg, lerr := localizer.Localize(&i18n.LocalizeConfig{
							MessageID:   em.code,
							PluralCount: em.capital,
							TemplateData: map[string]int{
								"digit": em.capital,
							},
						}); lerr == nil {
							msg = lmsg
						}
					default:
						log.Println("*checkPassErrors", "no val")
					}

					out = append(out, msg)
				}
				errFields[i] = out

			default:
				errFields[i] = []string{es.Error()}
			}
		}
	}

	return errFields
}

func (req *ChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.Login,
			is.Alphanumeric, validation.Required, validation.Length(2, 100)),
		validation.Field(&req.Password,
			validation.Required, validation.Length(5, 100), validation.By(checkSimplePass)),
	)
}

type checkPassErrors struct {
	errs []*checkPassError
}

func (e *checkPassErrors) Error() string {
	if len(e.errs) > 0 {
		errMsgs := make([]string, 0)
		for _, err := range e.errs {
			errMsgs = append(errMsgs, err.Error())
		}
		return strings.Join(errMsgs, "; ")
	}
	return ""
}

type checkPassError struct {
	digits, capital int
	code            string
}

func (e *checkPassError) Error() string {
	errMsgs := make([]string, 0)
	switch {
	case e.digits == 0 && e.capital == 0:
		errMsgs = append(errMsgs, e.code)

	case e.digits > 0:
		errMsgs = append(errMsgs, "there must be at least {{.digit}} digits in the value")

	case e.capital > 0:
		errMsgs = append(errMsgs, "there must be at least {{.capital}} capital letters in the value")
	}

	return strings.Join(errMsgs, "; ")
}

func checkSimplePass(value interface{}) error {
	var (
		s  string
		ok bool
	)
	if s, ok = value.(string); !ok {
		return &checkPassError{
			code: ErrNotString.Error(),
		}
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
		return &checkPassError{
			code: err.Error(),
		}
	}

	errs := make([]*checkPassError, 0)

	switch {
	case digits < _minDigit:
		errs = append(errs, &checkPassError{
			digits: _minDigit,
			code:   ValidationInternalErrNotEnoughDigits,
		})
		fallthrough

	case capital < _minCapital:
		errs = append(errs, &checkPassError{
			capital: _minCapital,
			code:    ValidationInternalErrNotEnoughCapitalLetters,
		})
	}

	if len(errs) > 0 {
		return &checkPassErrors{
			errs: errs,
		}
	}

	return nil
}
