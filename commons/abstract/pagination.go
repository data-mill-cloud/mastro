package abstract

import (
	paginate "github.com/gobeam/mongo-go-pagination"
)

type PaginationData struct {
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	PerPage   int64 `json:"perPage"`
	Prev      int64 `json:"prev"`
	Next      int64 `json:"next"`
	TotalPage int64 `json:"totalPage"`
}

type PaginatedAssets struct {
	Data       *[]Asset       `json:"data"`
	Pagination PaginationData `json:"pagination"`
}

func FromMongoPaginationData(pagination paginate.PaginationData) PaginationData {
	return PaginationData{
		Total:     pagination.Total,
		Page:      pagination.Page,
		PerPage:   pagination.PerPage,
		Prev:      pagination.Prev,
		Next:      pagination.Next,
		TotalPage: pagination.TotalPage,
	}
}
