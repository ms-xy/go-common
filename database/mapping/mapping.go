package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Mapping interface {
	GetType() reflect.Type
	GetFields() map[string]Field
	GetFieldNames() []string

	Scan(pRows SqlRows) (interface{}, error)
	MultiScan(pRows SqlRows, limit int) (interface{}, error)
}

/* -------------------------------------------------------------------------- */

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

/* -------------------------------------------------------------------------- */

type SqlRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

/* -------------------------------------------------------------------------- */

/*
ScanOnce scans a single row from the given `*sql.Rows` object.
The database cursor is NOT increased using `rows.Next()`, this has to be done separately.
Query fields are expected in mapping field order.
Returns a pointer to the newly created object of the mapped type.
*/
func (this *MappingStruct) Scan(pRows SqlRows) (interface{}, error) {
	pValue := reflect.New(this.Type)
	oValue := reflect.Indirect(pValue)
	aFields := make([]interface{}, len(this.FieldNames))
	for i := 0; i < len(this.FieldNames); i++ {
		oField := oValue.Field(i)
		aFields[i] = oField.Addr().Interface()
	}
	if err := pRows.Scan(aFields...); err == nil {
		return pValue.Interface(), nil
	} else {
		return nil, err
	}
}

// DefaultBatchScanSize specifies the default result array allocation size
// used with ScanLimit. Default value is 256.
var DefaultBatchScanSize = 256

/*
ScanLimit scans up to `limit` rows from the given `*sql.Rows` object,
advancing the cursor in doing so.
Upon encountering an error further scanning is aborted and the current result
plus error returned.
Return value is a slice of pointers to the mapped object.
Parameter `limit` defaults to `mapping.DefaultBachScanSize` (=256 default value).
A limit of -1 results in scanning all available rows.
Scanning with a limit of -1 is not adviced for large data sets, reflect.Copy is costly.
Instead consider scanning multiple times or using a large limit directly.
A limit of 0 results in returning an empty row-set.
*/
func (this *MappingStruct) MultiScan(pRows SqlRows, limit int) (interface{}, error) {
	if limit == 0 {
		return this.makeSlicePtr(0, 0).Elem().Interface(), nil
	} else {
		n := DefaultBatchScanSize
		if limit > 0 {
			n = limit
		}
		aResults := this.makeSlicePtr(n, n)
		i := 0
		for ; (limit == -1) || (i < limit); i++ {
			// next must only be called if limit has not been reached, thus can't be incorporated in
			// the loop
			ok := pRows.Next()
			if !ok {
				break
			}
			// old: for ok := pRows.Next(); ok && ((limit == -1) || (i < limit)); ok, i = pRows.Next(), i+1 {

			// if an error happened, return slice of results up until the scan
			if result, err := this.Scan(pRows); err != nil {
				aResults.Elem().SetLen(i)
				return aResults.Elem().Interface(), err
			} else {
				if i >= n {
					n = 2 * n
					newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(this.Type)), n, n)
					reflect.Copy(newSlice, aResults.Elem())
					aResults.Elem().Set(newSlice)
				}
				aResults.Elem().Index(i).Set(reflect.ValueOf(result))
			}
		}
		aResults.Elem().SetLen(i)
		return aResults.Elem().Interface(), nil
	}
}

func (this *MappingStruct) makeSlicePtr(len, cap int) reflect.Value {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(this.Type)), len, cap)
	slicePtr := reflect.New(slice.Type())
	slicePtr.Elem().Set(slice)
	return slicePtr
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
	ErrWrongType = errors.New("Supplied value not a struct")
)

/*
ValuesOf returns the values associated with the given structure in their defined
order, primarily for use with the transactions of the sibling package
*/
func ValuesOf(i interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(i)
	// indirect pointers until we have a struct
	if v.Kind() == reflect.Ptr {
		return ValuesOf(reflect.Indirect(v).Interface())
	}
	// require struct
	if v.Kind() != reflect.Struct {
		return nil, ErrWrongType
	}
	// iterate over all fields
	l := v.Type().NumField()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			return nil, fmt.Errorf(
				"Cannot use value.Interface() on struct field `%s`", v.Type().Field(i).Name)
		}
		r[i] = v.Field(i).Interface()
	}
	return r, nil
}

var (
	knownMappings     = make(map[reflect.Type]*MappingStruct)
	knownMappingsLock = sync.Mutex{}
	lock              = sync.Mutex{}
)

func ResetKnownMappings() {
	knownMappingsLock.Lock()
	defer knownMappingsLock.Unlock()
	knownMappings = make(map[reflect.Type]*MappingStruct)
}

/*
GetMapping returns a mapping instance for the given type.
Results are cached and subsequent requests for the same type are procured from
the internal cache.
*/
func GetMapping(i interface{}) (Mapping, error) {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		v := reflect.Indirect(reflect.ValueOf(i))
		t = v.Type()
		i = v.Interface()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrWrongType
	}
	// master lock to avoid race conditions when resetting mappings
	lock.Lock()
	defer lock.Unlock()
	// if a mapping exists, non needs to be created
	if mapping, exists := knownMappings[t]; exists {
		return mapping, nil
	} else {
		// lock mappings then check again, might have been created in between
		knownMappingsLock.Lock()
		defer knownMappingsLock.Unlock()
		if mapping, exists := knownMappings[t]; exists {
			return mapping, nil
		} else {
			// we have certainty that mapping does not exist and nobody is working on creating it
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
				fieldName := getFieldName(field)

				// write mapping
				mapping.FieldNames[i] = fieldName
				mapping.Fields[fieldName] = Field{
					Name:       field.Name,
					Properties: properties,
				}
			}

			knownMappings[t] = mapping
			return mapping, nil
		}
	}
}

var (
	fieldNameIdentifiers = []string{"db", "fieldname", "json"}
)

func SetFieldNameIdentifiers(identifiers []string) {
	lock.Lock()
	defer lock.Unlock()
	fieldNameIdentifiers = identifiers
	ResetKnownMappings()
}

func getFieldName(field reflect.StructField) string {
	for _, key := range fieldNameIdentifiers {
		if value, exists := field.Tag.Lookup(key); exists && (len(value) > 0) && (key != "json" || value != "-") {
			return value
		}
	}
	return field.Name
}
