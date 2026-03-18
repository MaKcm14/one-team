package user

import entity "github.com/MaKcm14/one-team/internal/entity/user"

type Claims struct {
	Login string      `json:"login"`
	Role  entity.Role `json:"role"`
}

type Credentials struct {
	Login    string
	Password string
}
