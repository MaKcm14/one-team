package entity

type DivisionID int

type DivisionType string

const (
	DivisionTypeName    DivisionType = "division"
	DirectorateTypeName DivisionType = "directorate"
	DepartmentTypeName  DivisionType = "department"
	UnitTypeName        DivisionType = "unit"
	GroupTypeName       DivisionType = "group"
)

type Division struct {
	Name            string
	Type            DivisionType
	StateSize       int
	SuperdivisionID DivisionID
}
