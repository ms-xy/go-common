package transactions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type TransactionHandler = func(ctx context.Context, tx *sql.Tx) (result interface{}, err error)

/*
WithReadTx is a convenience wrapper around WithTx, defining the sql-opts as
read-only.

See WithTx documentation for information about parameters and return value.
*/
func WithReadTx(db *sql.DB, ctx context.Context, fn TransactionHandler) (interface{}, error) {
	opts := new(sql.TxOptions)
	opts.Isolation = sql.LevelDefault
	opts.ReadOnly = true
	return WithTx(db, ctx, opts, fn)
}

/*
WithWriteTx is a convenience wrapper around WithTx.

See WithTx documentation for information about parameters and return value.
*/
func WithWriteTx(db *sql.DB, ctx context.Context, fn TransactionHandler) (interface{}, error) {
	opts := new(sql.TxOptions)
	opts.Isolation = sql.LevelDefault
	opts.ReadOnly = false
	return WithTx(db, ctx, opts, fn)
}

/*
WithTx starts a new transaction with ctx as parent context and opts specifying
the sql options to use. The given transaction handler function is executed in
the context of the transaction and provided with a child scoped context and
transaction handle.

If the supplied context is nil, an empty no-op context is automatically created.

Returns any errors and the result as provided by the transaction handler func.
*/
func WithTx(db *sql.DB, ctx context.Context, opts *sql.TxOptions, fn TransactionHandler) (interface{}, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		// in case of a panic, recover for rollback then throw again
		if r := recover(); r != nil {
			err := tx.Rollback()
			if _, ok := r.(error); !ok {
				r = errors.New(fmt.Sprintf("%s", r))
			}
			if err != nil {
				panic(errors.Wrap(r.(error), err.Error()))
			} else {
				panic(r)
			}
		}
	}()
	// run fn with sub context of given context for cancel cascading
	fnCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := fn(fnCtx, tx)
	// error within the handler rolls back as well
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			err = errors.Wrap(err, err2.Error())
		}
	} else {
		err = tx.Commit()
	}
	return result, err
}
