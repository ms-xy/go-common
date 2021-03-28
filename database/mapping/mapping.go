package mapping

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type Mapping interface {
	GetType() reflect.Type
	GetFields() map[string]Field
	GetFieldNames() []string

	ScanOnce(pRows *sql.Rows) (interface{}, error)
	ScanLimit(pRows *sql.Rows, limit int) ([]interface{}, error)
}

type MappingStruct struct {
	Type       reflect.Type
	Fields     map[string]Field
	FieldNames []string
}

func (this *MappingStruct) GetType() reflect.Type {
	return this.Type
}

func (this *MappingStruct) GetFields() map[string]Field {
	return this.Fields
}

func (this *MappingStruct) GetFieldNames() []string {
	return this.FieldNames
}

/*
ScanOnce scans a single row from the given `*sql.Rows` object, advancing the
cursor before scanning with `rows.Next()`.
Fields are expected in mapping field order.
*/
func (this *MappingStruct) ScanOnce(pRows *sql.Rows) (interface{}, error) {
	oValue := reflect.Indirect(reflect.New(this.Type))
	aFields := make([]interface{}, len(this.FieldNames))
	for i := 0; i < len(this.FieldNames); i++ {
		oField := oValue.Field(i)
		aFields[i] = oField.Addr().Interface()
	}
	if ok := pRows.Next(); ok {
		err := pRows.Scan(aFields...)
		return oValue.Interface(), err
	}
	return nil, pRows.Err()
}

// DefaultBatchScanSize specifies the default result array allocation size
// used with ScanLimit
var DefaultBatchScanSize = 256

/*
ScanLimit scans up to `limit` rows from the given `*sql.Rows` object,
advancing the cursor in doing so.
Upon encountering an error further scanning is aborted and the current result
plus error returned.
*/
func (this *MappingStruct) ScanLimit(pRows *sql.Rows, limit int) ([]interface{}, error) {
	i := 0
	aResults := make([]interface{}, 0)
	n := limit
	if limit <= 0 {
		n = DefaultBatchScanSize
	}
	for {
		j := 0
		aInterimsResults := make([]interface{}, n)
		for ok := pRows.Next(); ok; ok = pRows.Next() {
			if result, err := this.ScanOnce(pRows); err != nil {
				aInterimsResults[j] = err
			} else {
				aInterimsResults[j] = result
			}
			j++
			i++
		}
		aResults = append(aResults, aInterimsResults...)
		if limit > 0 || len(aInterimsResults) < n {
			break
		}
	}
	return aResults[:i], pRows.Err()
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
	return json.Marshal(struct {
		Type string
	}{
		Type: fp.Type.String(),
	})
}

/* -------------------------------------------------------------------------- */

var (
	ErrWrongType = errors.New("Can only map values of a struct type")
	ErrCantAddr  = errors.New("Inaccessible struct fields present")
)

/*
ValuesOf returns the values associated with the given structure in their defined
order, primarily for use with the transactions of the sibling package
*/
func ValuesOf(i interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return nil, ErrWrongType
	}
	l := v.Type().NumField()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		f := v.Field(i)
		if !f.CanAddr() {
			return nil, ErrCantAddr
		}
		r[i] = v.Field(i).Interface()
	}
	return r, nil
}

var (
	knownMappings = make(map[reflect.Type]*MappingStruct)
	lock          = sync.Mutex{}
)

/*
GetMapping returns a mapping instance for the given type.
Results are cached and subsequent requests for the same type are procured from
the internal cache.
*/
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
			mapping = &MappingStruct{
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
