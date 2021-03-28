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
UpdateObjects executes the given query len(objs) times.

UpdateObjects uses mapping.ValuesOf(obj) to get the args for each execution.
Any bound values associated with the query builder are ignored.

Return value numExecutions contains the value ob to the first error encountered,
if any.
Rows affected contains the total of rows affected.
*/
func UpdateObjects(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, objs []interface{}) (numExecutions int64, rowsAffected int64, err error) {
	if len(objs) < 1 {
		return 0, 0, ErrObjsEmpty
	}
	numExecutions = 0
	rowsAffected = 0
	err = WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		// ignore boundArgs
		queryString, _, err := query.ToSql()
		if err != nil {
			return err
		}
		for i := 0; i < len(objs); i++ {
			stmt, err := tx.PrepareContext(ctx, queryString)
			if err != nil {
				return err
			}
			args, err := mapping.ValuesOf(objs[i])
			if err != nil {
				return err
			}
			result, err := stmt.ExecContext(ctx, args...)
			affected, err := result.RowsAffected()
			rowsAffected += affected
			numExecutions++
		}
		return nil
	})
	return
}

/*
Update executes the given query.

If bound arguments are saved in the query builder, they are used.
*/
func Update(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder) (rowsAffected int64, err error) {
	rowsAffected = 0
	err = WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		queryString, boundArgs, err := query.ToSql()
		if err != nil {
			return err
		}
		stmt, err := tx.PrepareContext(ctx, queryString)
		if err != nil {
			return err
		}
		result, err := stmt.ExecContext(ctx, boundArgs...)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		rowsAffected = affected
		return err
	})
	return
}
