package mapping

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	io "github.com/ms-xy/go-common/io"
)

type Mapping struct {
	Type       reflect.Type
	Fields     map[string]Field
	FieldNames []string
}

// ScanOnce scans a single row from the given `*sql.Rows` object, it expects any
// error checks to happen before (you must call rows.Next() first yourself),
// further it expects all fields for the scan in the order of definition present
// in the struct used for the mapping.
func (this *Mapping) ScanOnce(pRows *sql.Rows, rowID int) (interface{}, error) {
	oValue := reflect.Indirect(reflect.New(this.Type))
	aFields := make([]interface{}, len(this.FieldNames))
	for i := 0; i < len(this.FieldNames); i++ {
		oField := oValue.Field(i)
		aFields[i] = oField.Addr().Interface()
	}
	if err := pRows.Scan(aFields...); err != nil {
		return nil, errors.New(
			fmt.Sprintf("Error scanning row #%d: %s", rowID, err.Error()))
	} else {
		return oValue.Interface(), nil
	}
}

// ScanLimit scans up to `limit` rows from the given `*sql.Rows` object,
// checking `rows.Next()` before every call as well as `rows.Err()`, an error
// object will be added to the result set if row could not be fetched.
func (this *Mapping) ScanLimit(pRows *sql.Rows, limit uint64) []interface{} {
	aResults := make([]interface{}, limit)
	i := 0
	for ok := pRows.Next(); ok || pRows.Err() != nil; ok = pRows.Next() {
		if err := pRows.Err(); err != nil {
			aResults[i] = errors.New(
				fmt.Sprintf("Error fetching row #%d: %s", i, err.Error()))
		} else {
			if result, err := this.ScanOnce(pRows, i); err != nil {
				aResults[i] = err
			} else {
				aResults[i] = result
			}
		}
		i++
	}
	return aResults[:i]
}

/* -------------------------------------------------------------------------- */

type Field struct {
	Name       string
	Properties FieldProperties
}

/* -------------------------------------------------------------------------- */

type FieldProperties struct {
	Type reflect.Type
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (fp FieldProperties) MarshalJSON() (bytes []byte, err error) {
	oWriter := io.NewByteWriter()
	oEncoder := json.NewEncoder(oWriter)
	oWriter.MustWrite([]byte(`{`))
	// oWriter.MustWrite([]byte(`"Alias":`))
	// must(oEncoder.Encode(fp.Alias))
	oWriter.MustWrite([]byte(`,"Type":`))
	must(oEncoder.Encode(fp.Type.String()))
	oWriter.MustWrite([]byte(`}`))
	return oWriter.Get()
}

/* -------------------------------------------------------------------------- */

var (
	knownMappings = make(map[reflect.Type]Mapping)
	lock          = sync.Mutex{}
)

func GetMapping(i interface{}) Mapping {
	t := reflect.TypeOf(i)
	if mapping, exists := knownMappings[t]; exists {
		return mapping
	} else {
		lock.Lock()
		defer lock.Unlock()
		if mapping, exists := knownMappings[t]; exists {
			return mapping
		} else {
			mapping = Mapping{
				Type:       t,
				Fields:     make(map[string]Field, t.NumField()),
				FieldNames: make([]string, t.NumField()),
			}

			for i, n := 0, t.NumField(); i < n; i++ {
				field := t.Field(i)

				// properties
				properties := FieldProperties{
					Type: field.Type,
				}

				// database field name
				var (
					exists    bool
					fieldName string
				)
				if fieldName, exists = field.Tag.Lookup("dbfield"); !exists {
					if fieldName, exists = field.Tag.Lookup("json"); !exists {
						fieldName = field.Name
					}
				}

				// write mapping
				mapping.FieldNames[i] = fieldName
				mapping.Fields[fieldName] = Field{
					Name:       field.Name,
					Properties: properties,
				}
			}

			knownMappings[t] = mapping
			return mapping
		}
	}
}
