package transactions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/ms-xy/go-common/database/mapping"
)

var (
	ErrObjsEmpty = errors.New("`objs` is empty")
)

/*
Insert executes the given query.

If bound arguments are saved in the query builder, they are used.
*/
func Insert(db *sql.DB, ctx context.Context, query squirrel.InsertBuilder) (sql.Result, error) {
	queryString, boundArgs, err := query.ToSql()
	if err != nil {
		return emptySqlResult(), err
	}
	return exec(db, ctx, queryString, boundArgs)
}

/*
InsertObjects executes the given query for every object provided.

If bound arguments are saved in the query builder, they are prepended to query
parameters.
*/
func InsertObjects(db *sql.DB, ctx context.Context, query squirrel.InsertBuilder, objs []interface{}) ([]sql.Result, error) {
	queryString, boundArgs, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return exec_objects(db, ctx, queryString, boundArgs, objs)
}

/*
SelectOne runs a given query with limit 1 agains the given database within a
read-only transaction, returning the mapped result or any error that occured.
*/
func SelectOne(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, m mapping.Mapping) (interface{}, error) {
	txResult, txErr := WithReadTx(db, ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		// enforce limit 1
		queryString, boundArgs, err := query.Limit(1).ToSql()

		stmt, err := tx.PrepareContext(ctx, queryString)
		if err != nil {
			return emptySqlResult(), err
		}

		rows, err := stmt.QueryContext(ctx, boundArgs...)
		if err != nil {
			return emptySqlResult(), err
		}
		defer rows.Close()

		return m.MultiScan(rows, 1)
	})
	return txResult, txErr
}

/*
Select runs a given query. If a `limit` is set on the given SelectBuilder
instance, only `limit` number of rows are fetched.
If an error is encounted, the error and all results up until the error are
returned.

Results are returned as an interface{} castable to []*<type>, example:

	type MyType struct {
		// ...
	}

	func main() {
		// ...
		r, err := transactions.Select(db, nil, query, mapping.GetMapping(MyType{}))
		// ... some err handling
		results := r.([]*MyType)
	}
*/
func Select(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, m mapping.Mapping) (interface{}, error) {
	// run select with read only transaction
	txResult, txErr := WithReadTx(db, ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		queryString, boundArgs, err := query.ToSql()

		stmt, err := tx.PrepareContext(ctx, queryString)
		if err != nil {
			return emptySqlResult(), err
		}

		rows, err := stmt.QueryContext(ctx, boundArgs...)
		if err != nil {
			return emptySqlResult(), err
		}
		defer rows.Close()

		// scan until no more rows
		return m.MultiScan(rows, -1)
	})
	return txResult, txErr
}

/*
Update executes the given query.

If bound arguments are saved in the query builder, they are used.
*/
func Update(db *sql.DB, ctx context.Context, query squirrel.UpdateBuilder) (sql.Result, error) {
	queryString, boundArgs, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return exec(db, ctx, queryString, boundArgs)
}

/*
UpdateObjects executes the given query len(objs) times.

UpdateObjects uses mapping.ValuesOf(obj) to get the args for each execution.
Any bound values associated with the query builder are prepended to object values.

Return value numExecutions contains the value ob to the first error encountered,
if any.
Rows affected contains the total of rows affected.
*/
func UpdateObjects(db *sql.DB, ctx context.Context, query squirrel.UpdateBuilder, objs []interface{}) ([]sql.Result, error) {
	if len(objs) < 1 {
		return nil, ErrObjsEmpty
	}
	queryString, boundArgs, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return exec_objects(db, ctx, queryString, boundArgs, objs)
}
