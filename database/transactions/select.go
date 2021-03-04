package transactions

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/ms-xy/go-common/database/mapping"
)

func GenericSelect(db *sql.DB, t interface{}, fromTable string, where interface{}, limit uint64) []interface{} {
	oMapping := mapping.GetMapping(t)
	oQuery := squirrel.Select(oMapping.FieldNames...).From(fromTable)

	if sWhere, ok := where.(string); ok {
		if where != "" {
			oQuery = oQuery.Where(sWhere)
		}
	} else if oWhere, ok := where.(squirrel.Sqlizer); ok {
		oQuery = oQuery.Where(oWhere)
	} else if where != nil {
		panic(errors.New("Unknown type for function parameter 'where': " + reflect.TypeOf(where).String()))
	}

	if limit > 0 {
		oQuery = oQuery.Limit(limit)
	}

	if pRows, err := oQuery.RunWith(db).Query(); err != nil {
		panic(err)

	} else {
		defer pRows.Close()
		return oMapping.ScanLimit(pRows, limit)
	}
}
