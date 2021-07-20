package transactions

import "database/sql"

type sqlResult struct {
	rowsAffected    int64
	lastInsertId    int64
	errRowsAffected error
	errLastInsertId error
}

func emptySqlResult() sql.Result {
	return sqlResult{}
}

func copySqlResult(r sql.Result) sql.Result {
	rowsAffected, errRowsAffected := r.RowsAffected()
	lastInsertId, errLastInsertId := r.LastInsertId()
	return sqlResult{
		rowsAffected:    rowsAffected,
		lastInsertId:    lastInsertId,
		errRowsAffected: errRowsAffected,
		errLastInsertId: errLastInsertId,
	}
}

func (s sqlResult) RowsAffected() (int64, error) {
	return s.rowsAffected, s.errRowsAffected
}

func (s sqlResult) LastInsertId() (int64, error) {
	return s.lastInsertId, s.errLastInsertId
}
