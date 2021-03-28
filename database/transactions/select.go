package transactions

import (
	"context"
	"database/sql"
	"math/bits"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/lann/builder"
	"github.com/ms-xy/go-common/database/mapping"
)

/*
SelectOne runs a given query with limit 1 agains the given database within a
read-only transaction, returning the result and/or any error that occured.
*/
func SelectOne(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, m mapping.Mapping) (interface{}, error) {
	var r interface{}
	err := WithReadTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		rows, err := query.Limit(1).RunWith(tx).QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()
		r, err = m.ScanOnce(rows)
		return err
	})
	return r, err
}

/*
Select runs a given query. If a `limit` is set on the given SelectBuilder
instance, only `limit` number of rows are fetched.
If an error is encounted, the error and all results up until the error are
returned.
*/
func Select(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, m mapping.Mapping) ([]interface{}, error) {
	// check if limit exists, if yes, parse it and apply
	limit := 0
	if_limit, ok := builder.Get(query, "Limit")
	if ok && if_limit != nil {
		// get the limit, platform dependent 32 or 64 bit
		_limit, err := strconv.ParseInt(if_limit.(string), 10, bits.UintSize)
		if err != nil {
			return nil, err
		}
		limit = int(_limit)
	}

	// run select with read only transaction
	var r []interface{} = make([]interface{}, 0)
	err := WithReadTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		rows, err := query.RunWith(tx).QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()

		if limit > 1 {
			// if `limit` applies (> 1) run `limit` number of times
			_r, _err := m.ScanLimit(rows, limit)
			r = _r
			err = _err
		} else {
			// if `limit` does not apply (< 2) only run once, saving a little
			// overhead
			_r, _err := m.ScanOnce(rows)
			r = make([]interface{}, 1)
			r[0] = _r
			err = _err
		}
		return err
	})
	return r, err
}
