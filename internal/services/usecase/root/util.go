package root

import entity "github.com/MaKcm14/one-team/internal/entity/user"

type UserDTO struct {
	Login string      `json:"login"`
	Role  entity.Role `json:"role"`
}

type Role struct {
	Name   string         `json:"name"`
	Rights []entity.Right `json:"rights"`
}
