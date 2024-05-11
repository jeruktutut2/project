package entity

import "database/sql"

type UserPermission struct {
	Id           sql.NullInt32
	UserId       sql.NullInt32
	PermissionId sql.NullInt32
}
