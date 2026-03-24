package entity

type DivisionType string

const (
	DivisionTypeName    DivisionType = "division"
	DirectorateTypeName DivisionType = "directorate"
	DepartmentTypeName  DivisionType = "department"
	UnitTypeName        DivisionType = "unit"
	GroupTypeName       DivisionType = "group"
)

func IsDivisionTypeValid(divType DivisionType) bool {
	switch divType {
	case DivisionTypeName, DirectorateTypeName, DepartmentTypeName, UnitTypeName, GroupTypeName:
		return true
	}
	return false
}

type Division struct {
	ID              int          `json:"division_id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Type            DivisionType `json:"type,omitempty"`
	StateSize       int          `json:"state_size,omitempty"`
	SuperdivisionID int          `json:"superdivision_id,omitempty"`
}
