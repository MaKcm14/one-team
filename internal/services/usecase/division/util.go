package division

import entity "github.com/MaKcm14/one-team/internal/entity/division"

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
