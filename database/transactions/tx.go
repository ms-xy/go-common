package transactions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type TransactionHandler = func(ctx context.Context, tx *sql.Tx) error

func WithReadTx(db *sql.DB, ctx context.Context, fn TransactionHandler) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	opts := new(sql.TxOptions)
	opts.Isolation = sql.LevelDefault
	opts.ReadOnly = true
	return WithTx(db, ctx, opts, fn)
}

func WithWriteTx(db *sql.DB, ctx context.Context, fn TransactionHandler) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	opts := new(sql.TxOptions)
	opts.Isolation = sql.LevelDefault
	opts.ReadOnly = false
	return WithTx(db, ctx, opts, fn)
}

func WithTx(db *sql.DB, ctx context.Context, opts *sql.TxOptions, fn TransactionHandler) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
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
			}
		}
	}()
	// run fn with sub context of given context for cancel cascading
	fnCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = fn(fnCtx, tx)
	// error within the handler rolls back as well
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}
