package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	ErrCodeUnknown ErrorCode = iota + 1000
	ErrCodeInternal
	ErrCodeBadRequest
	ErrCodeTooManyRequests
)
const (
	ErrCodeValidation ErrorCode = iota + 1100
	ErrCodeInvalidInput
	ErrCodeMissingField
)
const (
	ErrCodeDatabase ErrorCode = iota + 1200
	ErrCodeDuplicateEntry
)
const (
	ErrCodeNotFound ErrorCode = iota + 1300
	ErrCodeAlreadyExists
	ErrCodeConflict
)
const (
	ErrCodeUnauthorized ErrorCode = iota + 1400
	ErrCodeInvalidToken
	ErrCodeExpiredToken
	ErrCodeForbidden
	ErrCodeAuth
	ErrCodeAccountNotVerified
	ErrCodeInvalidProvider
	ErrCodeProviderTokenInvalid
	ErrCodeProviderEmailNotVerified
	ErrCodeAccountDeleted
)
const (
	ErrCodeExternalService ErrorCode = iota + 1500
	ErrCodeExternalServiceTimeout
	_
	ErrCodeTemplateNotFound
	ErrCodeTemplateRenderFailed
)

const (
	ErrCodeBusinessLogic ErrorCode = iota + 1600
)


// Error types
const (
	TypeUnknown       = "UNKNOWN_ERROR"
	TypeValidation    = "VALIDATION_ERROR"
	TypeNotFound      = "NOT_FOUND"
	TypeInternal      = "INTERNAL_ERROR"
	TypeUnauthorized  = "UNAUTHORIZED"
	TypeAuth          = "AUTH_ERROR"
	TypeForbidden     = "FORBIDDEN"
	TypeConflict      = "CONFLICT"
	TypeBadRequest    = "BAD_REQUEST"
	TypeDatabase      = "DATABASE_ERROR"
	TypeSocialAuth    = "SOCIAL_AUTH_ERROR"
	TypeExternal      = "EXTERNAL_SERVICE_ERROR"
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

func IsUnauthorizedError(err error) bool {
	return IsErrorCode(err, ErrCodeUnauthorized)
}

func InternalError(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, TypeInternal, message, err, http.StatusInternalServerError)
}

func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, TypeUnauthorized, message, err, http.StatusUnauthorized)
}

func ExpiredTokenError(message string, err error) *AppError {
	return NewAppError(ErrCodeExpiredToken, TypeUnauthorized, message, err, http.StatusUnauthorized)
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

func AccountNotVerifiedError(message string, err error) *AppError {
	return NewAppError(ErrCodeAccountNotVerified, "ACCOUNT_NOT_VERIFIED", message, err, http.StatusForbidden)
}

func AccountDeletedError(message string, err error) *AppError {
	return NewAppError(ErrCodeAccountDeleted, "ACCOUNT_DELETED", message, err, http.StatusForbidden)
}

func ExternalServiceError(message string, err error) *AppError {
	return NewAppError(ErrCodeExternalService, TypeExternal, message, err, http.StatusBadGateway)
}

// TooManyRequestsError creates a rate limiting error
func TooManyRequestsError(message string) *AppError {
	return NewAppError(ErrCodeTooManyRequests, TypeBadRequest, message, nil, http.StatusTooManyRequests)
}

// DatabaseError creates a generic database error
func DatabaseError(message string, err error) *AppError {
	return NewAppError(ErrCodeDatabase, TypeDatabase, message, err, http.StatusInternalServerError)
}

// DuplicateEntryError creates a duplicate entry error
func DuplicateEntryError(message string, err error) *AppError {
	return NewAppError(ErrCodeDuplicateEntry, TypeConflict, message, err, http.StatusConflict)
}

// InvalidInputError creates an invalid input error
func InvalidInputError(message string, err error) *AppError {
	return NewAppError(ErrCodeInvalidInput, TypeValidation, message, err, http.StatusBadRequest)
}

// MissingFieldError creates a missing field error
func MissingFieldError(fieldName string) *AppError {
	return NewAppError(ErrCodeMissingField, TypeValidation,
		fmt.Sprintf("Required field '%s' is missing", fieldName), nil, http.StatusBadRequest)
}

// AlreadyExistsError creates an already exists error
func AlreadyExistsError(message string, err error) *AppError {
	return NewAppError(ErrCodeAlreadyExists, TypeConflict, message, err, http.StatusConflict)
}

// InvalidTokenError creates an invalid token error
func InvalidTokenError(message string, err error) *AppError {
	return NewAppError(ErrCodeInvalidToken, TypeUnauthorized, message, err, http.StatusUnauthorized)
}

// ExternalServiceTimeoutError creates an external service timeout error
func ExternalServiceTimeoutError(message string, err error) *AppError {
	return NewAppError(ErrCodeExternalServiceTimeout, TypeExternal,
		message, err, http.StatusGatewayTimeout)
}

// BusinessLogicError creates a business logic error
func BusinessLogicError(message string, err error) *AppError {
	return NewAppError(ErrCodeBusinessLogic, "BUSINESS_LOGIC_ERROR", message, err, http.StatusBadRequest)
}

// SocialAuthError creates a social authentication error
func SocialAuthError(message string, err error) *AppError {
	return NewAppError(ErrCodeAuth, TypeSocialAuth, message, err, http.StatusUnauthorized)
}

// InvalidProviderError creates an invalid provider error
func InvalidProviderError(provider string) *AppError {
	return NewAppError(ErrCodeInvalidProvider, TypeSocialAuth,
		fmt.Sprintf("Invalid social auth provider: %s", provider), nil, http.StatusBadRequest)
}

// ProviderTokenError creates a provider token invalid error
func ProviderTokenError(message string, err error) *AppError {
	return NewAppError(ErrCodeProviderTokenInvalid, TypeSocialAuth, message, err, http.StatusUnauthorized)
}

// ProviderEmailNotVerifiedError creates a provider email not verified error
func ProviderEmailNotVerifiedError(provider string) *AppError {
	return NewAppError(ErrCodeProviderEmailNotVerified, TypeSocialAuth,
		fmt.Sprintf("Email not verified with %s provider", provider), nil, http.StatusForbidden)
}

func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
