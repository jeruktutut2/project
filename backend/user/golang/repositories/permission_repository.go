package repository

import (
	"context"
	"database/sql"
	modelentities "project-user/models/entities"
)

type PermissionRepository interface {
	FindByInId(db *sql.DB, ctx context.Context, ids []interface{}) (permissions []modelentities.Permission, err error)
}

type PermissinoRepositoryImplementation struct {
}

func NewPermissionRepository() PermissionRepository {
	return &PermissinoRepositoryImplementation{}
}

func (repository *PermissinoRepositoryImplementation) FindByInId(db *sql.DB, ctx context.Context, ids []interface{}) (permissions []modelentities.Permission, err error) {
	var placeholder string
	for i := 0; i < len(ids); i++ {
		placeholder += ",?"
	}
	placeholder = placeholder[1:]
	rows, err := db.QueryContext(ctx, `SELECT id, permission FROM permission IN (`+placeholder+`)`, ids...)
	if err != nil {
		return
	}
	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			permissions = []modelentities.Permission{}
			err = errRowsClose
		}
	}()

	for rows.Next() {
		var permission modelentities.Permission
		err = rows.Scan(&permission.Id, &permission.Permission)
		if err != nil {
			permissions = []modelentities.Permission{}
			return
		}
		permissions = append(permissions, permission)
	}
	return
}
