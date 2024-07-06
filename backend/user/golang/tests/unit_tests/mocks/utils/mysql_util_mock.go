package mockutils

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MysqlUtilMock struct {
	Mock mock.Mock
}

func (util *MysqlUtilMock) GetDb() *sql.DB {
	arguments := util.Mock.Called()
	return arguments.Get(0).(*sql.DB)
}

func (util *MysqlUtilMock) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	arguments := util.Mock.Called(ctx, options)
	return arguments.Get(0).(*sql.Tx), arguments.Error(1)
}

func (util *MysqlUtilMock) Close() {
	arguments := util.Mock.Called()
	fmt.Println(arguments)
}

func (util *MysqlUtilMock) CommitOrRollback(tx *sql.Tx, err error) error {
	arguments := util.Mock.Called(tx, err)
	return arguments.Error(0)
}
