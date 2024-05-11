package entity

import "database/sql"

type Permission struct {
	Id         sql.NullInt32
	Permission sql.NullString
}
