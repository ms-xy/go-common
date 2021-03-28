package dbtypes

import (
	"database/sql/driver"
	"encoding/json"

	convert "github.com/Eun/go-convert"
)

/*
	NullString represents a string that may be null.

	NullString implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer

	Example Usage:

		var s NullString
		err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
		...
		if !s.IsNull() {
		  // use s.Get()
		} else {
		  // NULL value
		}
*/
type NullString struct {
	str  string
	null bool
}

func (n NullString) Set(v *string) {
	if v == nil {
		n.str = ""
		n.null = true
	} else {
		n.str = (*v)[:]
		n.null = false
	}
}

func (n NullString) Get() string {
	return n.str
}

func (n NullString) IsNull() bool {
	return n.null
}

func (this *NullString) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v string
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullString) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullString) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullString) UnmarshalJSON(bytes []byte) error {
	var v *string
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullInt64 represents an int64 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullInt64 struct {
	value int64
	null  bool
}

func (this NullInt64) Set(v *int64) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullInt64) Get() int64 {
	return this.value
}

func (this NullInt64) IsNull() bool {
	return this.null
}

func (this *NullInt64) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v int64
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullInt64) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullInt64) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullInt64) UnmarshalJSON(bytes []byte) error {
	var v *int64
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullInt32 represents an int64 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullInt32 struct {
	value int32
	null  bool
}

func (this NullInt32) Set(v *int32) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullInt32) Get() int32 {
	return this.value
}

func (this NullInt32) IsNull() bool {
	return this.null
}

func (this *NullInt32) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v int32
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullInt32) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullInt32) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullInt32) UnmarshalJSON(bytes []byte) error {
	var v *int32
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullInt16 represents an int16 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullInt16 struct {
	value int16
	null  bool
}

func (this NullInt16) Set(v *int16) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullInt16) Get() int16 {
	return this.value
}

func (this NullInt16) IsNull() bool {
	return this.null
}

func (this *NullInt16) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v int16
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullInt16) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullInt16) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullInt16) UnmarshalJSON(bytes []byte) error {
	var v *int16
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullInt8 represents an int64 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullInt8 struct {
	value int8
	null  bool // Valid is true if Int32 is not NULL
}

func (this NullInt8) Set(v *int8) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullInt8) Get() int8 {
	return this.value
}

func (this NullInt8) IsNull() bool {
	return this.null
}

func (this *NullInt8) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v int8
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullInt8) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullInt8) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullInt8) UnmarshalJSON(bytes []byte) error {
	var v *int8
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullFloat64 represents a float64 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullFloat64 struct {
	value float64
	null  bool // Valid is true if Float64 is not NULL
}

func (this NullFloat64) Set(v *float64) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullFloat64) Get() float64 {
	return this.value
}

func (this NullFloat64) IsNull() bool {
	return this.null
}

func (this *NullFloat64) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v float64
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullFloat64) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullFloat64) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullFloat64) UnmarshalJSON(bytes []byte) error {
	var v *float64
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}

/*
	NullFloat64 represents a float64 that may be null.

	NullInt64 implements
	- encoding/json.Mashaller
	- encoding/json.Unmarshaller
	- sql.Scanner
	- sql/driver.Valuer
*/
type NullFloat32 struct {
	value float32
	null  bool
}

func (this NullFloat32) Set(v *float32) {
	if v == nil {
		this.value = 0
		this.null = true
	} else {
		this.value = *v
		this.null = false
	}
}

func (this NullFloat32) Get() float32 {
	return this.value
}

func (this NullFloat32) IsNull() bool {
	return this.null
}

func (this *NullFloat32) Scan(value interface{}) error {
	var err error
	if value != nil {
		var v float32
		err = convert.Convert(value, &v)
		if err == nil {
			this.Set(&v)
			return nil
		}
	}
	this.Set(nil)
	return err
}

func (this NullFloat32) Value() (driver.Value, error) {
	if this.IsNull() {
		return nil, nil
	}
	return this.Get(), nil
}

func (this NullFloat32) MarshalJSON() ([]byte, error) {
	if this.IsNull() {
		var v interface{} = nil
		return json.Marshal(v)
	} else {
		return json.Marshal(this.Get())
	}
}

func (this NullFloat32) UnmarshalJSON(bytes []byte) error {
	var v *float32
	err := json.Unmarshal(bytes, &v)
	if err == nil {
		this.Set(v)
	} else {
		this.Set(nil)
	}
	return err
}
