package mockdb

import (
  "context"
  "database/sql/driver"

  "github.com/stretchr/testify/mock"
)

/* implement driver.Conn */

type Conn struct {
  mock.Mock
}

func (this *Conn) Prepare(query string) (driver.Stmt, error) {
  r := this.Called(query)
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Stmt), err
}

func (this *Conn) Close() error {
  r := this.Called()
  err, _ := r.Get(0).(error)
  return err
}

func (this *Conn) Begin() (driver.Tx, error) {
  r := this.Called()
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Tx), err
}

func (this *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
  r := this.Called(ctx, opts)
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Tx), err
}
