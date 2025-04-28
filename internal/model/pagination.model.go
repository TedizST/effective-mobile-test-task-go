package model

import (
	"effective-mobile-test-task/internal/types"
)

type Pagination struct {
	Limit types.Limit
	Page  types.Page
}

func (p Pagination) GetLimit() uint64 {
	if p.Limit == 0 {
		return 10
	}
	return uint64(p.Limit)
}

func (p Pagination) GetPage() uint64 {
	if p.Page == 0 {
		return 1
	}
	return uint64(p.Page)
}
