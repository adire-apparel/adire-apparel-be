package dto

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"items"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func Paginate[T any](data []T, params Pagination) PaginatedResponse[T] {
	totalItems := int64(len(data))
	totalPages := int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))

	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start > len(data) {
		start = len(data)
	}

	if end > len(data) {
		end = len(data)
	}

	return PaginatedResponse[T]{
		Data:       data[start:end],
		Limit:      int(params.Limit),
		Page:       params.Page,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}
}
