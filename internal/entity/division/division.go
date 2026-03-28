package entity

type DivisionType string

const (
	DivisionTypeName    DivisionType = "division"
	DirectorateTypeName DivisionType = "directorate"
	DepartmentTypeName  DivisionType = "department"
	UnitTypeName        DivisionType = "unit"
	GroupTypeName       DivisionType = "group"
	None                DivisionType = "none"
)

func IsDivisionTypeValid(divType DivisionType) bool {
	switch divType {
	case DivisionTypeName, DirectorateTypeName, DepartmentTypeName, UnitTypeName, GroupTypeName:
		return true
	}
	return false
}

// TODO: consider to replace this checking in the separate checking in migrations (in CHECK f.e).
func IsDivisionTypeRelationCorrect(div DivisionType, supDiv DivisionType) bool {
	switch div {
	case DivisionTypeName:
		if supDiv != None {
			return false
		}
		return true

	case DirectorateTypeName:
		if supDiv != DivisionTypeName {
			return false
		}
		return true

	case DepartmentTypeName:
		if supDiv != DirectorateTypeName && supDiv != DivisionTypeName {
			return false
		}
		return true

	case UnitTypeName:
		if supDiv != DirectorateTypeName && supDiv != DivisionTypeName && supDiv != DepartmentTypeName {
			return false
		}
		return true

	case GroupTypeName:
		if supDiv == GroupTypeName {
			return false
		}
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
