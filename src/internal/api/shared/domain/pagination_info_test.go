//+build !e2e

package domain_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/stretchr/testify/require"
)

func TestNewPaginationInfo(t *testing.T) {
	testCases := []struct {
		limit               int
		offset              int
		sortField           string
		order               domain.PaginationOrder
		expectedStringOrder string
	}{
		{10, 20, "fieldName", domain.OrderAsc, "fieldName asc"},
		{10, 20, "fieldName", domain.OrderDesc, "fieldName desc"},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("New pagination info for %v '%v'", c.sortField, c.order), func(t *testing.T) {
			info := domain.NewPaginationInfo(c.limit, c.offset, c.sortField, c.order)

			require.NotNil(t, info)
			require.IsType(t, &domain.PaginationInfo{}, info)
			require.Equal(t, c.limit, info.Limit)
			require.Equal(t, c.offset, info.Offset)
			require.Equal(t, c.expectedStringOrder, info.Order)
		})
	}
}

func TestNewPaginationInfoFromUrl(t *testing.T) {
	testCases := []struct {
		url                 string
		expectedLimit       int
		expectedOffset      int
		expectedStringOrder string
	}{
		{"/url", 10, 0, "id asc"},
		{"/url?page=0&page_size=20", 20, 0, "id asc"},
		{"/url?page=2&page_size=-10", 10, 10, "id asc"},
		{"/url?page=2&page_size=200", 100, 100, "id asc"},
		{"/url?page=2&page_size=30&sort=fieldName", 30, 30, "fieldName asc"},
		{"/url?page=2&page_size=30&sort=fieldName&order=asc", 30, 30, "fieldName asc"},
		{"/url?page=2&page_size=30&sort=fieldName&order=desc", 30, 30, "fieldName desc"},
		{"/url?page=2&page_size=30&sort=fieldName&order=wadus", 30, 30, "fieldName asc"},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("New pagination info for %v", c.url), func(t *testing.T) {
			u, _ := url.Parse(c.url)
			info := domain.NewPaginationInfoFromUrl(u)

			require.NotNil(t, info)
			require.IsType(t, &domain.PaginationInfo{}, info)
			require.Equal(t, c.expectedLimit, info.Limit, "wrong limit")
			require.Equal(t, c.expectedOffset, info.Offset, "wrong offset")
			require.Equal(t, c.expectedStringOrder, info.Order, "wrong string order")
		})
	}
}
