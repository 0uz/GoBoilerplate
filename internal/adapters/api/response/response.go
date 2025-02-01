package response

import (
	"encoding/json"
	"net/http"

	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitResponseLogger(l *logrus.Logger) {
	logger = l
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
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
