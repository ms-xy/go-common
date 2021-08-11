package mockdb

import (
  "database/sql/driver"

  "github.com/stretchr/testify/mock"
)

/* implement the sql-driver interface */

type Driver struct {
  mock.Mock
}

func (this *Driver) Open(name string) (driver.Conn, error) {
  r := this.Called(name)
  err, _ := r.Get(1).(error)
  return r.Get(0).(*Conn), err
}
