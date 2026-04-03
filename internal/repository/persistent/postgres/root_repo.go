package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/root"
)

type rootRepo struct {
	client *postgresClient
}

const getRoleRightsQuery = `
SELECT app_realm.rights.name
FROM
	app_realm.roles
	JOIN
	app_realm.role_rights_mapping
	ON app_realm.role_rights_mapping.role_id=app_realm.roles.id
	JOIN
	app_realm.rights
	ON app_realm.role_rights_mapping.right_id=app_realm.rights.id
WHERE app_realm.roles.name=$1;
`

func (r rootRepo) getRoleRights(ctx context.Context, roleName entity.Role) ([]entity.Right, error) {
	res, err := r.client.conn.Query(ctx, getRoleRightsQuery, roleName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]entity.Right, 0, 100)
	for res.Next() {
		var right entity.Right
		if err := res.Scan(&right); err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		list = append(list, right)
	}

	return list, nil
}

const getRolesQuery = `
SELECT name
FROM app_realm.roles;
`

func (r rootRepo) GetRoles(ctx context.Context) ([]root.Role, error) {
	res, err := r.client.conn.Query(ctx, getRolesQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]root.Role, 0, 100)
	for res.Next() {
		var role root.Role

		if err := res.Scan(&role.Name); err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}

		rights, err := r.getRoleRights(ctx, entity.Role(role.Name))
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		role.Rights = rights

		list = append(list, role)
	}
	return list, nil
}
