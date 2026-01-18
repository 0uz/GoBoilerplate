package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ouz/goboilerplate/pkg/errors"
	"github.com/ouz/goboilerplate/pkg/log"
)

var Validator = validator.New()

var logger *log.Logger

func InitResponseLogger(l *log.Logger) {
	logger = l
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode response", "error", err)
	}
}

func Error(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		JSON(w, appErr.Status, appErr)
		return
	}

	JSON(w, http.StatusInternalServerError, errors.InternalError("An unexpected error occurred", err))
}

func DecodeAndValidate(r *http.Request, request any) error {
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := Validator.Struct(request); err != nil {
		return errors.BadRequestError(err.Error())
	}

	return nil
}

func CreatePagination[T any](r *http.Request) (*Pagination[T], error) {
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page, limit := 1, 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			return nil, errors.BadRequestError("Page must be a positive integer")
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l >= 1 && l <= 100 {
			limit = l
		} else {
			return nil, errors.BadRequestError("Limit must be between 1 and 100")
		}
	}

	return &Pagination[T]{
		Page:  page,
		Limit: limit,
	}, nil
}
