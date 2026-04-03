package user

import (
	"context"
)

type IAuthService interface {
	Login(ctx context.Context, creds Credentials) (UserDTO, error)
	SignUp(ctx context.Context, dto UserSignUpDTO) error
	ChangePassword(ctx context.Context, creds Credentials, newPwd string) error
}
