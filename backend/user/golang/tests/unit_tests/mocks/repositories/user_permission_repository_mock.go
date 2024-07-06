package mockrepositories

import (
	"context"
	"database/sql"

	modelentities "project-user/models/entities"

	"github.com/stretchr/testify/mock"
)

type UserPermissionRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserPermissionRepositoryMock) FindByUserId(db *sql.DB, ctx context.Context, userId int32) (userPermissions []modelentities.UserPermission, err error) {
	arguments := repository.Mock.Called(db, ctx, userId)
	return arguments.Get(0).([]modelentities.UserPermission), arguments.Error(1)
}
