package datastructures

import "github.com/google/uuid"

// Set is a very simple interface definition for minimal set datastructure
// functionality.
// Will be a lot easier and better to use once go has generics ...
type Set interface {
	// Contains checks if the given value is contained within the set
	Contains(interface{}) bool
	// Add adds the given value to the set
	Add(interface{})
	// Remove removes the given value from the set
	Remove(interface{})
	// Slice returns the values contained within the set as a slice
	Slice() interface{}
}

type UuidSet struct {
	m map[uuid.UUID]struct{}
}

func NewUuidSet() *UuidSet {
	return &UuidSet{m: make(map[uuid.UUID]struct{})}
}

func (this *UuidSet) Contains(i interface{}) bool {
	_, exists := this.m[i.(uuid.UUID)]
	return exists
}

func (this *UuidSet) Add(i interface{}) {
	this.m[i.(uuid.UUID)] = struct{}{}
}

func (this *UuidSet) Remove(i interface{}) {
	delete(this.m, i.(uuid.UUID))
}

func (this *UuidSet) Slice() interface{} {
	r := make([]uuid.UUID, len(this.m))
	i := 0
	for u := range this.m {
		r[i] = u
	}
	return r
}

type StringSet struct {
	m map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{m: make(map[string]struct{})}
}

func (this *StringSet) Contains(i interface{}) bool {
	_, exists := this.m[i.(string)]
	return exists
}

func (this *StringSet) Add(i interface{}) {
	this.m[i.(string)] = struct{}{}
}

func (this *StringSet) Remove(i interface{}) {
	delete(this.m, i.(string))
}

func (this *StringSet) Slice() interface{} {
	r := make([]string, len(this.m))
	i := 0
	for u := range this.m {
		r[i] = u
	}
	return r
}
