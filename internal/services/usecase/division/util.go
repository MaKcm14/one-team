package division

import (
	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

const (
	// FilterTypes.
	NameFilterName = "name"

	PaginationSize = 25
)

type NameFilter struct {
	Name    string
	Type    entity.DivisionType
	PageNum int
}

type Filters struct {
	Names NameFilter
}

type SalaryStatistics struct {
	Average float64 `json:"avg"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

type StateSizeStatistics struct {
	MinStateSizeDivList []entity.Division `json:"min_state_size_divisions"`
	MaxStateSizeDivList []entity.Division `json:"max_state_size_divisions"`
}
