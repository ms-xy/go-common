package slices

import (
	"reflect"
)

// Filter applies fn on each value in the input slice and returns a new slice
// containing only those values where fn yielded true
func Filter(slice interface{}, fn func(reflect.Value) bool) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("cannot use Filter with datatypes other than slice or array")
	}
	l := v.Len()
	if l == 0 {
		return slice
	}

	r := reflect.MakeSlice(v.Type(), l, l)
	for i, j := 0, 0; i < l; i++ {
		f := v.Index(i)
		if fn(f) {
			r.Index(j).Set(f)
		}
	}
	return r
}
