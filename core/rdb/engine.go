package rdb

import (
	"context"
)

type Engine[T any] interface {
	FetchOne(ctx context.Context, table string, sql string, args ...interface{}) (*T, error)
	FetchAll(ctx context.Context, table interface{}, sql string, args ...interface{}) ([]*T, error)
	Exec(ctx context.Context, table, sql string, args ...interface{}) (int64, error)
	Count(ctx context.Context, table, sql string, args ...interface{}) (int64, error)
}
