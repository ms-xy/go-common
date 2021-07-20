package transactions

import (
	"context"
	"database/sql"

	"github.com/ms-xy/go-common/database/mapping"
)

/*
execute runs a query against a database.
*/
func exec(db *sql.DB, ctx context.Context, queryString string, boundArgs []interface{}) (sql.Result, error) {
	txResult, txErr := WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		stmt, err := tx.PrepareContext(ctx, queryString)
		if err != nil {
			return emptySqlResult(), err
		}
		if boundArgs == nil {
			boundArgs = make([]interface{}, 0)
		}
		result, err := stmt.ExecContext(ctx, boundArgs...)
		if err != nil {
			return emptySqlResult(), err
		}
		return copySqlResult(result), nil
	})
	return txResult.(sql.Result), txErr
}

/*
exec_objects runs a query against a database for each object supplied.

mapping.ValuesOf(obj) is used internally to get the args for each execution.
Any bound values associated with the query builder are prepended to this list of
arguments in the call to tx.ExecContext(ctx, args...).

Return value is a slice of type []sql.Result up until the first error encountered or until
all queries are executed.
*/
func exec_objects(db *sql.DB, ctx context.Context, queryString string, boundArgs []interface{}, objs []interface{}) ([]sql.Result, error) {
	if len(objs) < 1 {
		return nil, ErrObjsEmpty
	}

	txResult, txErr := WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		i := 0
		sqlResults := make([]sql.Result, len(objs))

		for ; i < len(objs); i++ {
			stmt, _err := tx.PrepareContext(ctx, queryString)
			if _err != nil {
				return sqlResults[:i], _err
			}
			args, _err := mapping.ValuesOf(objs[i])
			if _err != nil {
				return sqlResults[:i], _err
			}
			if len(boundArgs) != 0 {
				args = append(boundArgs, args...)
			}
			result, _err := stmt.ExecContext(ctx, args...)
			if _err != nil {
				return sqlResults[:i], _err
			}
			sqlResults[i] = copySqlResult(result)
		}

		return sqlResults[:i], nil
	})

	return txResult.([]sql.Result), txErr
}
