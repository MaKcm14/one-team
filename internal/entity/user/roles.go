package entity

const (
	AdminRole     Role = "admin"
	HRManagerRole Role = "hr-manager"
	AnalystRole   Role = "analyst"
)

type (
	Right string
	Role  string
)
