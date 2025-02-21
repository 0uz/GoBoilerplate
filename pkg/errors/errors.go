package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	// General errors
	ErrCodeUnknown ErrorCode = iota + 1000
	ErrCodeInternal
	ErrCodeBadRequest
	ErrCodeTooManyRequests

	// Validation errors
	ErrCodeValidation ErrorCode = iota + 2000
	ErrCodeInvalidInput
	ErrCodeMissingField

	// Authentication errors
	ErrCodeUnauthorized ErrorCode = iota + 3000
	ErrCodeForbidden
	ErrCodeInvalidToken
	ErrCodeAuth

	// Resource errors
	ErrCodeNotFound ErrorCode = iota + 4000
	ErrCodeAlreadyExists
	ErrCodeConflict

	// Database errors
	ErrCodeDatabase ErrorCode = iota + 5000
	ErrCodeDuplicateEntry

	// External service errors
	ErrCodeExternalService ErrorCode = iota + 6000
	ErrCodeExternalServiceTimeout

	// Business logic errors
	ErrCodeBusinessLogic ErrorCode = iota + 7000
)

// Error types
const (
	TypeUnknown      = "UNKNOWN_ERROR"
	TypeValidation   = "VALIDATION_ERROR"
	TypeNotFound     = "NOT_FOUND"
	TypeInternal     = "INTERNAL_ERROR"
	TypeUnauthorized = "UNAUTHORIZED"
	TypeAuth         = "AUTH_ERROR"
	TypeForbidden    = "FORBIDDEN"
	TypeConflict     = "CONFLICT"
	TypeBadRequest   = "BAD_REQUEST"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
	Status  int       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Error %d: %s - %s: %v", e.Code, e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("Error %d: %s - %s", e.Code, e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
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

func IsErrorCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

func IsErrorType(err error, errType string) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == errType
	}
	return false
}

func GenericError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnknown, TypeUnknown, message, err, http.StatusInternalServerError)
}

func ValidationError(message string, err error) *AppError {
	return NewAppError(ErrCodeValidation, TypeValidation, message, err, http.StatusBadRequest)
}

func NotFoundError(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, TypeNotFound, message, err, http.StatusNotFound)
}

func IsNotFoundError(err error) bool {
	return IsErrorCode(err, ErrCodeNotFound)
}

func InternalError(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, TypeInternal, message, err, http.StatusInternalServerError)
}

func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, TypeUnauthorized, message, err, http.StatusUnauthorized)
}

func AuthError(message string, err error) *AppError {
	return NewAppError(ErrCodeAuth, TypeAuth, message, err, http.StatusUnauthorized)
}

func ForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrCodeForbidden, TypeForbidden, message, err, http.StatusForbidden)
}

func ConflictError(message string, err error) *AppError {
	return NewAppError(ErrCodeConflict, TypeConflict, message, err, http.StatusConflict)
}

func BadRequestError(message string) *AppError {
	return NewAppError(ErrCodeBadRequest, TypeBadRequest, message, nil, http.StatusBadRequest)
}

func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func (e *AppError) GetStatus() int {
	return e.Status
}
