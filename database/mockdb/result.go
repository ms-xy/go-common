package mockdb

import "github.com/stretchr/testify/mock"

type Result struct {
  mock.Mock
}

func (this *Result) LastInsertId() (int64, error) {
  r := this.Called()
  err, _ := r.Get(1).(error)
  return r.Get(0).(int64), err
}

func (this *Result) RowsAffected() (int64, error) {
  r := this.Called()
  err, _ := r.Get(1).(error)
  return r.Get(0).(int64), err
}
