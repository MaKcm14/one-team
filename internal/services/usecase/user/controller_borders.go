package user

import (
	"context"
)

type IAuthService interface {
	Login(ctx context.Context, creds Credentials) error
}
