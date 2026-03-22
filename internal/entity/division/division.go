package entity

type DivisionType string

const (
	DivisionTypeName    DivisionType = "division"
	DirectorateTypeName DivisionType = "directorate"
	DepartmentTypeName  DivisionType = "department"
	UnitTypeName        DivisionType = "unit"
	GroupTypeName       DivisionType = "group"
)

type Division struct {
	ID              int          `json:"division_id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Type            DivisionType `json:"type,omitempty"`
	StateSize       int          `json:"state_size,omitempty"`
	SuperdivisionID int          `json:"superdivision_id,omitempty"`
}
