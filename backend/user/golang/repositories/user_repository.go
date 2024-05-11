package repository

import (
	"context"
	"database/sql"
	modelentity "project-user/models/entities"
)

type UserRepository interface {
	Create(db *sql.DB, ctx context.Context, user modelentity.User) (rowsAffected int64, err error)
	CountByUsername(db *sql.DB, ctx context.Context, username string) (numberOfUser int, err error)
	CountByEmail(db *sql.DB, ctx context.Context, email string) (numbeOfUser int, err error)
	FindByEmail(db *sql.DB, ctx context.Context, email string) (user modelentity.User, err error)
}

type UserRepositoryImplementation struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImplementation{}
}

func (repository *UserRepositoryImplementation) Create(db *sql.DB, ctx context.Context, user modelentity.User) (rowsAffected int64, err error) {
	result, err := db.ExecContext(ctx, `INSERT INTO user (username, email, password, utc, created_at) VALUES(?, ?, ?, ?, ?);`, user.Username, user.Email, user.Password, user.Utc, user.CreatedAt)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *UserRepositoryImplementation) CountByUsername(db *sql.DB, ctx context.Context, username string) (numberOfUser int, err error) {
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) AS number_of_user FROM user WHERE username = ?;`, username).Scan(&numberOfUser)
	return
}

func (repository *UserRepositoryImplementation) CountByEmail(db *sql.DB, ctx context.Context, email string) (numbeOfUser int, err error) {
	err = db.QueryRowContext(ctx, `SELECT COUNT(*) AS number_of_user FROM user WHERE email = ?;`, email).Scan(&numbeOfUser)
	return
}

func (repository *UserRepositoryImplementation) FindByEmail(db *sql.DB, ctx context.Context, email string) (user modelentity.User, err error) {
	err = db.QueryRowContext(ctx, `SELECT id, username, email, password, utc, created_at FROM user WHERE email = ?;`, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Utc, &user.CreatedAt)
	return
}
