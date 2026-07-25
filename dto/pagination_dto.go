package dto

// PaginationQuery di-bind dari query params: ?page=1&limit=10
type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

// Normalize menjaga dari nilai yang tidak wajar atau kosong, alih-alih
// mempercayai klien untuk selalu mengirim nilai yang benar.
func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

func (p PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResponse membungkus payload list apapun dengan metadata paginasi.
type PaginatedResponse struct {
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func BuildPaginatedResponse(data any, p PaginationQuery, totalItems int64) PaginatedResponse {
	totalPages := int(totalItems) / p.Limit
	if int(totalItems)%p.Limit != 0 {
		totalPages++
	}
	return PaginatedResponse{
		Data:       data,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
