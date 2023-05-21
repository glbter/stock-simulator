package sql

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"
)

type Selector interface {
	Select(ctx context.Context, dst any, query string, args ...any) error
	SelectOne(ctx context.Context, dst any, query string, args ...any) error
}

type Executor interface {
	Insert(ctx context.Context, id any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type SelectExecutor interface {
	Selector
	Executor
}

type DB interface {
	SelectExecutor

	Close() error
	SetSchema(ctx context.Context) error
	Migrate(ctx context.Context, fs embed.FS, dir, table string) (int, error)
}

var ErrNotFound = errors.New("no rows")

func IsSafeSQL(q string) bool {
	return !strings.Contains(q, ",\";':|*")

}

func PrepareSchema(ctx context.Context, exec Executor, schema string) error {
	if _, err := exec.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)); err != nil {
		return err
	}

	return nil
}
