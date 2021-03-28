package transactions

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
)

/*
Delete executes the given query.

If bound arguments are saved in the query builder, they are used.

Currently this is simply a duplicate of Update.
*/
func Delete(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder) (rowsAffected int64, err error) {
	rowsAffected = 0
	err = WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		queryString, boundArgs, _err := query.ToSql()
		if _err != nil {
			return _err
		}
		stmt, _err := tx.PrepareContext(ctx, queryString)
		if _err != nil {
			return _err
		}
		result, _err := stmt.ExecContext(ctx, boundArgs...)
		if _err != nil {
			return _err
		}
		affected, _err := result.RowsAffected()
		if _err != nil {
			return _err
		}
		rowsAffected = affected
		return nil
	})
	return
}
