package user

import entity "github.com/MaKcm14/one-team/internal/entity/user"

type UserDTO struct {
	User entity.User
	Role entity.Role
}

type UserSignUpDTO struct {
	Creds Credentials `json:"credentials"`
	Role  entity.Role `json:"role_name"`
}
