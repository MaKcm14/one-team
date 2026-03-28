package employee

import (
	division "github.com/MaKcm14/one-team/internal/entity/division"
	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

const (
	// FilterTypes.
	PassportDataFilterName = "passport_data"
	NamesFilterName        = "names"
	UnitFilterName         = "unit"
	UnitTypeFilterName     = "unit_type"

	PaginationSize = 25
)

type PassportFilter struct {
	IsActive     bool
	PassportData string
	PageNum      int
}

type NamesFilter struct {
	IsActive   bool
	FirstName  string
	LastName   string
	Patronymic string
	PageNum    int
}

type UnitFilter struct {
	IsActive bool
	Name     string
	Type     division.DivisionType
	PageNum  int
}

type Filter struct {
	Passport PassportFilter
	Names    NamesFilter
	Unit     UnitFilter
}

type EmployeeCitizenshipStatistic struct {
	Citizenship    entity.Citizenship `json:"citizenship"`
	EmployeesCount int                `json:"employees_count"`
}

type SalaryBounds struct {
	UpBoundary   int
	DownBoundary int
}
