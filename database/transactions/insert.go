package transactions

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/ms-xy/go-common/database/mapping"
)

var (
//ErrObjsEmpty = errors.New("`objs` is empty")
)

/*
InsertObjects executes the given query len(objs) times.

InsertObjects uses mapping.ValuesOf(obj) to get the args for each execution.
Any bound values associated with the query builder are ignored.

Return value numExecutions contains the value ob to the first error encountered,
if any.
Rows affected contains the total of rows affected.
*/
func InsertObjects(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder, objs []interface{}) (numExecutions int64, insertIds []int64, rowsAffected int64, err error) {
	if len(objs) < 1 {
		return 0, nil, 0, ErrObjsEmpty
	}
	numExecutions = 0
	insertIds = make([]int64, len(objs))
	rowsAffected = 0
	err = WithWriteTx(db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		// ignore boundArgs
		queryString, _, err := query.ToSql()
		if err != nil {
			return err
		}
		for i := 0; i < len(objs); i++ {
			stmt, _err := tx.PrepareContext(ctx, queryString)
			if _err != nil {
				return _err
			}
			args, _err := mapping.ValuesOf(objs[i])
			if _err != nil {
				return _err
			}
			result, _err := stmt.ExecContext(ctx, args...)
			if _err != nil {
				return _err
			}
			affected, _err := result.RowsAffected()
			if _err != nil {
				return _err
			}
			lastInserted, _err := result.LastInsertId()
			if _err != nil {
				return _err
			}
			insertIds[i] = lastInserted
			rowsAffected += affected
			numExecutions++
		}
		return nil
	})
	return
}

/*
Insert executes the given query.

If bound arguments are saved in the query builder, they are used.
*/
func Insert(db *sql.DB, ctx context.Context, query squirrel.SelectBuilder) (insertId, rowsAffected int64, err error) {
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
		lastInsertId, _err := result.LastInsertId()
		if _err != nil {
			return _err
		}
		insertId = lastInsertId
		rowsAffected = affected
		return _err
	})
	return
}
