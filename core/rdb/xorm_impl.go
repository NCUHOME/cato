package rdb

import (
	"context"
	"errors"

	"xorm.io/xorm"
)

type sessionKey struct{}

func NewXormEngine[T any](provider func(ctx context.Context) *xorm.Engine) Engine[T] {
	return &XormEngine[T]{provider}
}

type XormEngine[T any] struct {
	provider func(ctx context.Context) *xorm.Engine
}

func (engine *XormEngine[T]) session(ctx context.Context) *xorm.Session {
	if session, ok := ctx.Value(sessionKey{}).(*xorm.Session); ok {
		return session
	}
	return engine.provider(ctx).Context(ctx)
}

func (engine *XormEngine[T]) FetchOne(ctx context.Context, table string, sql string, args ...interface{}) (*T, error) {
	data := new(T)
	ok, err := engine.session(ctx).Table(table).SQL(sql, args...).Get(data)
	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return data, nil
}

func (engine *XormEngine[T]) FetchAll(ctx context.Context, table string, sql string, args ...interface{}) ([]*T, error) {
	data := make([]*T, 0)
	err := engine.session(ctx).Table(table).SQL(sql, args...).Find(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (engine *XormEngine[T]) Exec(ctx context.Context, table, sql string, args ...interface{}) (int64, error) {
	packs := make([]interface{}, len(args)+1)
	packs[0] = sql
	packs = append(packs, args...)
	result, err := engine.session(ctx).Table(table).Exec(packs...)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (engine *XormEngine[T]) Count(ctx context.Context, table, sql string, args ...interface{}) (int64, error) {
	result, err := engine.session(ctx).Table(table).Where(sql, args...).Count()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func NewXormTx(ctx context.Context, sessionCreate func(ctx context.Context) *xorm.Session) (context.Context, Tx) {
	sess := sessionCreate(ctx)
	return context.WithValue(ctx, sessionKey{}, sess), &XormTx{sess: sess}
}

type XormTx struct {
	sess *xorm.Session
}

func (tx *XormTx) Begin() error {
	return tx.sess.Begin()
}

func (tx *XormTx) Commit() error {
	return errors.Join(tx.sess.Commit(), tx.sess.Close())
}

func (tx *XormTx) Rollback() error {
	return errors.Join(tx.sess.Rollback(), tx.sess.Close())
}
