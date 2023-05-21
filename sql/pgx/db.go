package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	sqlc "github.com/glbter/currency-ex/sql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type result struct {
	pgconn.CommandTag
}

func (r result) RowsAffected() (int64, error) {
	return r.CommandTag.RowsAffected(), nil
}

func (r result) LastInsertId() (int64, error) {
	panic("LastInsertID is not supported in PostgreSQL")
}

var _ sqlc.DB = &DB{}

type DB struct {
	conpool     *Pool
	schema      string
	schemaQuery string
}

func (db *DB) IsInit() bool {
	return db.conpool != nil
}

func NewDB(pool *Pool, schema string) *DB {
	const public = "public"

	var path []string
	if schema == "" || schema == public {
		path = []string{public}
	} else {
		path = []string{schema}
		//path = []string{schema, public}

	}

	return &DB{
		conpool:     pool,
		schema:      path[0],
		schemaQuery: fmt.Sprintf("SET search_path TO %s", strings.Join(path, ",")),
	}
}

func (db *DB) setSchema(ctx context.Context, conn *pgxpool.Conn) error {
	if _, err := conn.Exec(ctx, db.schemaQuery); err != nil {
		return err
	}

	return nil
}

func (db *DB) SetSchema(ctx context.Context) error {
	_, err := db.Exec(ctx, db.schemaQuery)
	return err
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := db.conpool.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, convertError(err)
	}

	return result{res}, nil
}

func (db *DB) Insert(ctx context.Context, id any, query string, args ...any) error {
	if err := db.conpool.pool.QueryRow(ctx, query, args...).Scan(id); err != nil {
		return convertError(err)
	}

	return nil
}

func (db *DB) Select(ctx context.Context, dst any, query string, args ...any) error {
	if err := pgxscan.Select(ctx, db.conpool.pool, dst, query, args...); err != nil {
		return convertError(err)
	}

	return nil
}

func (db *DB) SelectOne(ctx context.Context, dst any, query string, args ...any) error {
	if err := pgxscan.Get(ctx, db.conpool.pool, dst, query, args...); err != nil {
		return convertError(err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.conpool.Close()
}

func convertError(e error) error {
	if pgxscan.NotFound(e) {
		return fmt.Errorf("%w: %v", sqlc.ErrNotFound, e)
	}

	return e
}
