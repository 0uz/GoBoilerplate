package response

type Pagination[T any] struct {
	Limit      int
	Page       int
	TotalRows  int64 `json:"totalRows"`
	TotalPages int   `json:"totalPages"`
	Data       []T   `json:"data"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"totalRows"`
	TotalPages int   `json:"totalPages"`
	Data       any   `json:"data"`
}

func (p *Pagination[T]) GetLimit() int {
	if p.Limit <= 0 {
		return 10 // Default limit
	}
	if p.Limit > 100 {
		return 100 // Maximum limit
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page <= 0 {
		return 1 // Default page
	}
	return p.Page
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func ToPaginatedResponse[T any](pagination *Pagination[T], data any) *PaginationResponse {
	return &PaginationResponse{
		Page:       pagination.GetPage(),
		Limit:      pagination.GetLimit(),
		TotalRows:  pagination.TotalRows,
		TotalPages: pagination.TotalPages,
		Data:       data,
	}
}
