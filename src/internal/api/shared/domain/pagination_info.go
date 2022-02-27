package domain

import (
	"fmt"
	"net/url"
	"strconv"
)

type PaginationInfo struct {
	Limit  int
	Offset int
	Order  string
}

func NewPaginationInfo(limit int, offset int, sortField string, order PaginationOrder) *PaginationInfo {
	return &PaginationInfo{limit, offset, fmt.Sprintf("%v %v", sortField, order.String())}
}

func NewPaginationInfoFromUrl(url *url.URL) *PaginationInfo {
	page, _ := strconv.Atoi(url.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(url.Query().Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	sort := url.Query().Get("sort")
	order := url.Query().Get("order")

	if sort == "" {
		sort = "id"
	}

	if order == "" {
		order = "asc"
	}

	paginationOrder := OrderDesc

	if order != "desc" {
		paginationOrder = OrderAsc
	}

	return NewPaginationInfo(pageSize, offset, sort, paginationOrder)
}
