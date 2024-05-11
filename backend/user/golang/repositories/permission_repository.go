package repository

import (
	"context"
	"database/sql"
	modelentity "project-user/models/entities"
)

type PermissionRepository interface {
	FindByInId(db *sql.DB, ctx context.Context, ids []interface{}) (permissions []modelentity.Permission, err error)
}

type PermissinoRepositoryImplementation struct {
}

func NewPermissionRepository() PermissionRepository {
	return &PermissinoRepositoryImplementation{}
}

func (repository *PermissinoRepositoryImplementation) FindByInId(db *sql.DB, ctx context.Context, ids []interface{}) (permissions []modelentity.Permission, err error) {
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
			permissions = []modelentity.Permission{}
			err = errRowsClose
		}
	}()

	for rows.Next() {
		var permission modelentity.Permission
		err = rows.Scan(&permission.Id, &permission.Permission)
		if err != nil {
			permissions = []modelentity.Permission{}
			return
		}
		permissions = append(permissions, permission)
	}
	return
}
