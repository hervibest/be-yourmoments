package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Querier interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}
