package dbtypes

import (
	"database/sql/driver"
	"strconv"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var s NullString
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  if s.Valid {
//     // use s.String
//  } else {
//     // NULL value
//  }
//
// NullString implements the encoding.json.Mashaller interface to properly
// support json null
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (this *NullString) Scan(value interface{}) error {
	if value == nil {
		this.String, this.Valid = "", false
		return nil
	}
	this.Valid = true
	return convertAssign(&this.String, value)
}

// Value implements the driver Valuer interface.
func (this NullString) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.String, nil
}

func (this NullString) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(`"` + this.String + `"`), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullString) Set(value string) {
	this.String = value
	this.Valid = true
}

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
// NullInt64 implements the encoding.json.Mashaller interface to properly
// support json null
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.
func (this *NullInt64) Scan(value interface{}) error {
	if value == nil {
		this.Int64, this.Valid = 0, false
		return nil
	}
	this.Valid = true
	return convertAssign(&this.Int64, value)
}

// Value implements the driver Valuer interface.
func (this NullInt64) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.Int64, nil
}

func (this NullInt64) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(strconv.FormatInt(this.Int64, 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullInt64) Set(value int64) {
	this.Int64 = value
	this.Valid = true
}

// NullInt32 represents an int64 that may be null.
// NullInt32 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
// NullInt32 implements the encoding.json.Mashaller interface to properly
// support json null
type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt32) Scan(value interface{}) error {
	if value == nil {
		n.Int32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int32, value)
}

// Value implements the driver Valuer interface.
func (this NullInt32) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.Int32, nil
}

func (this NullInt32) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(strconv.FormatInt(int64(this.Int32), 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullInt32) Set(value int32) {
	this.Int32 = value
	this.Valid = true
}

// NullInt8 represents an int64 that may be null.
// NullInt8 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
// NullInt8 implements the encoding.json.Mashaller interface to properly
// support json null
type NullInt8 struct {
	Int8  int8
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt8) Scan(value interface{}) error {
	if value == nil {
		n.Int8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int8, value)
}

// Value implements the driver Valuer interface.
func (this NullInt8) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.Int8, nil
}

func (this NullInt8) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(strconv.FormatInt(int64(this.Int8), 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullInt8) Set(value int8) {
	this.Int8 = value
	this.Valid = true
}

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface.
func (this *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		this.Float64, this.Valid = 0, false
		return nil
	}
	this.Valid = true
	return convertAssign(&this.Float64, value)
}

// Value implements the driver Valuer interface.
func (this NullFloat64) Value() (driver.Value, error) {
	if !this.Valid {
		return nil, nil
	}
	return this.Float64, nil
}

func (this NullFloat64) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(strconv.FormatFloat(this.Float64, 'f', -1, 64)), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullFloat64) Set(value float64) {
	this.Float64 = value
	this.Valid = true
}

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Float32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Float32, value)
}

// Value implements the driver Valuer interface.
func (n NullFloat32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float32, nil
}

func (this NullFloat32) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return []byte(strconv.FormatFloat(float64(this.Float32), 'f', -1, 64)), nil
	} else {
		return []byte("null"), nil
	}
}

func (this NullFloat32) Set(value float32) {
	this.Float32 = value
	this.Valid = true
}
