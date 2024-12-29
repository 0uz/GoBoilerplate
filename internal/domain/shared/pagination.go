package shared

type Pagination[T any] struct {
	Limit      int
	Page       int
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Data       []T   `json:"data"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
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
