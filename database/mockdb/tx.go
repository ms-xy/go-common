package mockdb

import "github.com/stretchr/testify/mock"

type Tx struct {
  mock.Mock
}

func (this *Tx) Commit() error {
  r := this.Called()
  err, _ := r.Get(0).(error)
  return err
}

func (this *Tx) Rollback() error {
  r := this.Called()
  err, _ := r.Get(0).(error)
  return err
}
