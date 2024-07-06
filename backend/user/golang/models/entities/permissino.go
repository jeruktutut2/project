package entities

import "database/sql"

type Permission struct {
	Id         sql.NullInt32
	Permission sql.NullString
}
