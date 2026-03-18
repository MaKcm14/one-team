package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
)

type userRepo struct {
	client *postgresClient
}

const getUserQuery = `
	SELECT login, hash_pwd, salt
	FROM app_realm.users
	WHERE login=$1
`

func (u userRepo) GetUser(ctx context.Context, login string) (entity.User, error) {
	res, err := u.client.conn.Query(ctx, getUserQuery, login)
	if err != nil {
		retErr := fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		u.client.log.Error(res.Err().Error())
		return entity.User{}, retErr
	}
	defer res.Close()

	if !res.Next() {
		return entity.User{}, persistent.ErrUserNotFound
	}

	user := entity.User{}
	if err := res.Scan(&user.Login, &user.HashPWD, &user.Salt); err != nil {
		retErr := fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		u.client.log.Error(res.Err().Error())
		return entity.User{}, retErr
	}
	return user, nil
}

const getUserRoleQuery = `
SELECT app_realm.roles.name
FROM app_realm.users
	JOIN 
	app_realm.roles
	ON app_realm.users.role_id=app_realm.roles.id
WHERE app_realm.users.login=$1
`

func (u userRepo) GetUserRole(ctx context.Context, login string) (entity.Role, error) {
	res, err := u.client.conn.Query(ctx, getUserRoleQuery, login)
	if err != nil {
		return "", fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if res.Next() {
		var role entity.Role
		if err := res.Scan(&role); err != nil {
			return "", fmt.Errorf("%w: %s", persistent.ErrQueryExec)
		}
		return role, nil
	}
	return "", persistent.ErrRoleNotAssign
}
