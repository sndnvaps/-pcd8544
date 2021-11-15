package main

import (
	"github.com/cheekybits/genny/generic"
)

type Key generic.Type
type Value generic.Type

type ValueDictionary struct {
	data map[Key][5]Value
}

func NewValueDictionary() *ValueDictionary {
	return &ValueDictionary{
		data: map[Key][5]Value{},
	}
}
func (s *ValueDictionary) Set(key Key, value [5]Value) {
	if s.data == nil {
		s.data = map[Key][5]Value{}
	}
	s.data[key] = value
}

func (s *ValueDictionary) Delete(key Key) bool {
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

func (s *ValueDictionary) Has(key Key) bool {
	_, result := s.data[key]

	return result
}

func (s *ValueDictionary) Get(key Key) [5]Value {
	result, _ := s.data[key]
	return result
}

func (s *ValueDictionary) Clear() {
	s.data = map[Key][5]Value{}
}

func (s *ValueDictionary) Size() int {
	return len(s.data)
}

func (s *ValueDictionary) Keys() []Key {
	keys := make([]Key, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *ValueDictionary) Values() [][5]Value {
	values := make([][5]Value, len(s.data))
	for _, v := range s.data {
		values = append(values, v)
	}
	return values
}
