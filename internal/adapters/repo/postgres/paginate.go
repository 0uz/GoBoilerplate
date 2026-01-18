package postgres

import (
	"math"
	"net/http"
	"strconv"

	"github.com/ouz/goauthboilerplate/internal/domain/shared"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"gorm.io/gorm"
)

func Paginate[T any](value any, pagination *shared.Pagination[T], db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func CreatePagination[T any](r *http.Request) (*shared.Pagination[T], error) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}
	if limit > 100 || limit < 1 {
		return nil, errors.BadRequestError("Limit must be between 1 and 100")
	}
	if page < 1 {
		return nil, errors.BadRequestError("Page must be greater than 0")
	}
	return &shared.Pagination[T]{
		Page:  page,
		Limit: limit,
	}, nil
}

func CacheKeyWithPagination[T any](pagination *shared.Pagination[T]) string {
	if pagination == nil {
		return ""
	}
	return strconv.Itoa(pagination.GetPage()) + ":" + strconv.Itoa(pagination.GetLimit())
}
