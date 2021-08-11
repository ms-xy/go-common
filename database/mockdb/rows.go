package mockdb

import (
  "database/sql/driver"

  "github.com/stretchr/testify/mock"
)

type Rows struct {
  mock.Mock
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice. If a particular column name isn't known, an empty
// string should be returned for that entry.
func (this *Rows) Columns() []string {
  r := this.Called()
  return r.Get(0).([]string)
}

// Close closes the rows iterator.
func (this *Rows) Close() error {
  r := this.Called()
  err, _ := r.Get(0).(error)
  return err
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (this *Rows) Next(dest []driver.Value) error {
  r := this.Called(dest)
  err, _ := r.Get(0).(error)
  return err
}
