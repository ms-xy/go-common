package mockdb

import (
  "database/sql/driver"

  "github.com/stretchr/testify/mock"
)

type Stmt struct {
  mock.Mock
}

func (this *Stmt) Close() error {
  r := this.Called()
  err, _ := r.Get(0).(error)
  return err
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (this *Stmt) NumInput() int {
  r := this.Called()
  return r.Get(0).(int)
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (this *Stmt) Exec(args []driver.Value) (driver.Result, error) {
  r := this.Called(args)
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Result), err
}

// Query executes a query that may return rows, such as a
// SELECT.
//
// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (this *Stmt) Query(args []driver.Value) (driver.Rows, error) {
  r := this.Called(args)
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Rows), err
}
