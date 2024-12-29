package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	// Genel hata kodları
	ErrCodeUnknown    ErrorCode = 1000
	ErrCodeInternal   ErrorCode = 1001
	ErrCodeBadRequest ErrorCode = 1002

	// Doğrulama hata kodları
	ErrCodeValidation   ErrorCode = 2000
	ErrCodeInvalidInput ErrorCode = 2001
	ErrCodeMissingField ErrorCode = 2002

	// Yetkilendirme hata kodları
	ErrCodeUnauthorized ErrorCode = 3000
	ErrCodeForbidden    ErrorCode = 3001
	ErrCodeInvalidToken ErrorCode = 3002
	ErrCodeAuth         ErrorCode = 3003

	// Kaynak hata kodları
	ErrCodeNotFound      ErrorCode = 4000
	ErrCodeAlreadyExists ErrorCode = 4001
	ErrCodeConflict      ErrorCode = 4002

	// Veritabanı hata kodları
	ErrCodeDatabase       ErrorCode = 5000
	ErrCodeDuplicateEntry ErrorCode = 5001

	// Dış servis hata kodları
	ErrCodeExternalService        ErrorCode = 6000
	ErrCodeExternalServiceTimeout ErrorCode = 6001

	// İş mantığı hata kodları
	ErrCodeBusinessLogic ErrorCode = 7000
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
	Status  int       `json:"-"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("Error %d: %s - %s", e.Code, e.Type, e.Message)
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code ErrorCode, errType string, message string, err error, status int) *AppError {
	return &AppError{
		Code:    code,
		Type:    errType,
		Message: message,
		Err:     err,
		Status:  status,
	}
}

func GenericError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnknown, "UNKNOWN_ERROR", message, err, http.StatusInternalServerError)
}

func ValidationError(message string, err error) *AppError {
	return NewAppError(ErrCodeValidation, "VALIDATION_ERROR", message, err, http.StatusBadRequest)
}

func NotFoundError(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, "NOT_FOUND", message, err, http.StatusNotFound)
}

func IsNotFoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

func InternalError(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, "INTERNAL_ERROR", message, err, http.StatusInternalServerError)
}

func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, "UNAUTHORIZED", message, err, http.StatusUnauthorized)
}

func AuthError(message string, err error) *AppError {
	return NewAppError(ErrCodeAuth, "AUTH_ERROR", message, err, http.StatusUnauthorized)
}

func ForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrCodeForbidden, "FORBIDDEN", message, err, http.StatusForbidden)
}

func ConflictError(message string, err error) *AppError {
	return NewAppError(ErrCodeConflict, "CONFLICT", message, err, http.StatusConflict)
}

func BadRequestError(message string) *AppError {
	return NewAppError(ErrCodeBadRequest, "BAD_REQUEST", message, nil, http.StatusBadRequest)
}
