package mockrepositories

import (
	"context"
	"database/sql"

	modelentities "project-user/models/entities"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) Create(db *sql.DB, ctx context.Context, user modelentities.User) (rowsAffected int64, err error) {
	arguments := repository.Mock.Called(db, ctx, user)
	return arguments.Get(0).(int64), arguments.Error(1)
}

func (repository *UserRepositoryMock) CountByUsername(db *sql.DB, ctx context.Context, username string) (numberOfUser int, err error) {
	arguments := repository.Mock.Called(db, ctx, username)
	return arguments.Get(0).(int), arguments.Error(1)
}

func (repository *UserRepositoryMock) CountByEmail(db *sql.DB, ctx context.Context, email string) (numbeOfUser int, err error) {
	arguments := repository.Mock.Called(db, ctx, email)
	return arguments.Get(0).(int), arguments.Error(1)
}

func (repository *UserRepositoryMock) FindByEmail(db *sql.DB, ctx context.Context, email string) (user modelentities.User, err error) {
	arguments := repository.Mock.Called(db, ctx, email)
	return arguments.Get(0).(modelentities.User), arguments.Error(1)
}
