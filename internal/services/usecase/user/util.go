package user

import entity "github.com/MaKcm14/one-team/internal/entity/user"

const (
	PaginationSize = 25
)

type Claims struct {
	Login string      `json:"login"`
	Role  entity.Role `json:"role"`
}

type UserSession struct {
	UserClaims Claims
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginFilter struct {
	IsActive bool
	Login    string
	PageNum  int
}

type Filters struct {
	LoginFilter LoginFilter
}
