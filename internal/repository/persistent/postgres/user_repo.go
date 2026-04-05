package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/root"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
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

const getRoleIDByNameQuery = `
SELECT id
FROM app_realm.roles
WHERE name=$1;
`

func (u userRepo) getRoleIDByName(ctx context.Context, name entity.Role) (int, error) {
	res, err := u.client.conn.Query(ctx, getRoleIDByNameQuery, name)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return 0, persistent.ErrRoleNotFound
	}

	var id int
	if err := res.Scan(&id); err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return id, nil
}

const createUserQuery = `
INSERT INTO app_realm.users (login, hash_pwd, salt, role_id)
VALUES ($1, $2, $3, $4);
`

func (u userRepo) CreateUser(ctx context.Context, dto user.UserDTO) error {
	id, err := u.getRoleIDByName(ctx, dto.Role)
	if err != nil {
		return err
	}

	_, err = u.client.conn.Exec(ctx, createUserQuery, dto.User.Login, dto.User.HashPWD, dto.User.Salt, id)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return nil
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
			return "", fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		return role, nil
	}
	return "", persistent.ErrRoleNotAssign
}

const getUsersQuery = `
SELECT app_realm.users.login, app_realm.users.hash_pwd, app_realm.users.salt, app_realm.roles.name
FROM 
	app_realm.users
	JOIN
	app_realm.roles
	ON app_realm.users.role_id=app_realm.roles.id;
`

func (u userRepo) GetUsers(ctx context.Context) ([]user.UserDTO, error) {
	res, err := u.client.conn.Query(ctx, getUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]user.UserDTO, 0, 100)
	for res.Next() {
		var user user.UserDTO
		err := res.Scan(&user.User.Login, &user.User.HashPWD, &user.User.Salt, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		list = append(list, user)
	}

	return list, nil
}

const getUsersByLoginQuery = `
SELECT app_realm.users.login, app_realm.users.hash_pwd, app_realm.users.salt, app_realm.roles.name
FROM 
	app_realm.users
	JOIN
	app_realm.roles
	ON app_realm.users.role_id=app_realm.roles.id
WHERE login LIKE $1
OFFSET $2
LIMIT $3;
`

func (u userRepo) GetUsersByLogin(ctx context.Context, filter user.LoginFilter) ([]user.UserDTO, error) {
	res, err := u.client.conn.Query(
		ctx,
		getUsersByLoginQuery,
		as(filter.Login),
		filter.PageNum*user.PaginationSize,
		user.PaginationSize,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]user.UserDTO, 0, 100)
	for res.Next() {
		var user user.UserDTO
		err := res.Scan(&user.User.Login, &user.User.HashPWD, &user.User.Salt, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		list = append(list, user)
	}

	return list, nil
}

const deleteUserQuery = `
DELETE FROM app_realm.users
WHERE login=$1;
`

func (u userRepo) DeleteUser(ctx context.Context, login string) error {
	res, err := u.client.conn.Exec(ctx, deleteUserQuery, login)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrUserNotFound
	}
	return nil
}

const updateUserRoleQuery = `
UPDATE app_realm.users
SET role_id=(
	SELECT id
	FROM app_realm.roles
	WHERE name=$1
)
WHERE login=$2;
`

func (u userRepo) UpdateUserRole(ctx context.Context, user root.UserDTO) error {
	res, err := u.client.conn.Exec(
		ctx,
		updateUserRoleQuery,
		user.Role,
		user.Login,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrUserNotFound
	}
	return nil
}

const updateUserPasswordQuery = `
UPDATE app_realm.users
SET salt=$1, hash_pwd=$2
WHERE login=$3;
`

func (u userRepo) UpdateUserPassword(ctx context.Context, user entity.User) error {
	res, err := u.client.conn.Exec(
		ctx,
		updateUserPasswordQuery,
		user.Salt,
		user.HashPWD,
		user.Login,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrUserNotFound
	}
	return nil
}
