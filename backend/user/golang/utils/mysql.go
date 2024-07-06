package utils

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// why do this, because when do unit test, there is an error that said something like invalid memory pointer etc, and also if use mock library, when do begintx, i dont know, just appear error that said everything is ...

type MysqlUtil interface {
	GetDb() *sql.DB
	BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error)
	Close()
	CommitOrRollback(tx *sql.Tx, err error) error
}

type MysqlUtilImplementation struct {
	db *sql.DB
}

func NewMysqlConnection(username string, password string, host string, database string, maxOpenconnection int, maxIdleConnection int, connectionMaxLifetime int, connectionMaxIdletime int) MysqlUtil {
	println(time.Now().String(), "mysql: connecting to", host)

	// why use concat instead fmt.Sprintf, because there is warning leaking param with username, password, host, database
	// dsn := username + `:` + password + `@tcp(` + host + `:` + port + `)/` + database + `?charset=utf8mb4&parseTime=True&loc=Local`
	dsn := username + `:` + password + `@tcp(` + host + `)/` + database + `?charset=utf8mb4&parseTime=True&loc=Local`
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenconnection)
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetConnMaxLifetime(time.Minute * time.Duration(connectionMaxLifetime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(connectionMaxIdletime))

	println(time.Now().String(), "mysql: connected to", host)
	return &MysqlUtilImplementation{
		db: db,
	}
}

func (util *MysqlUtilImplementation) GetDb() *sql.DB {
	return util.db
}

func (util *MysqlUtilImplementation) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	return util.db.BeginTx(ctx, options)
}

func (util *MysqlUtilImplementation) Close() {
	err := util.db.Close()
	if err != nil {
		panic(err)
	}
	println(time.Now().String(), "mysql: closed properly")
}

func (util *MysqlUtilImplementation) CommitOrRollback(tx *sql.Tx, err error) error {
	if err == nil {
		errCommit := tx.Commit()
		if errCommit != nil && !errors.Is(errCommit, sql.ErrTxDone) {
			errRollback := tx.Rollback()
			if errRollback != nil && !errors.Is(errRollback, sql.ErrTxDone) {
				return errRollback
			}
			return nil
		}
		return nil
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil && !errors.Is(errRollback, sql.ErrTxDone) {
			return errRollback
		}
		return nil
	}
}
