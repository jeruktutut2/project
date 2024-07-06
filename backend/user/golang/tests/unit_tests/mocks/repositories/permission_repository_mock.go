package mockrepositories

import (
	"context"
	"database/sql"
	modelentities "project-user/models/entities"

	"github.com/stretchr/testify/mock"
)

type PermissionRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PermissionRepositoryMock) FindByInId(db *sql.DB, ctx context.Context, ids []interface{}) (permissions []modelentities.Permission, err error) {
	arguments := repository.Mock.Called(db, ctx, ids)
	return arguments.Get(0).([]modelentities.Permission), arguments.Error(1)
}
