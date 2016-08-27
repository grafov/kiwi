package kiwi

import (
	"strconv"
)

type filter interface {
	Check(string, string) bool
}

type keyFilter struct {
	Key string
}

func (k *keyFilter) Check(key, val string) bool {
	return k.Key == val
}

type valsFilter struct {
	Key  string
	Vals []string
}

func (k *valsFilter) Check(key, val string) bool {
	if key != k.Key {
		return false
	}
	for _, v := range k.Vals {
		if v == val {
			return true
		}
	}
	return false
}

type rangeInt64Filter struct {
	Key      string
	From, To int64
}

func (k *rangeInt64Filter) Check(key, val string) bool {
	if key != k.Key {
		return false
	}
	var (
		intVal int64
		err    error
	)
	if intVal, err = strconv.ParseInt(val, 10, 64); err != nil {
		return false
	}
	return intVal > k.From && intVal <= k.To
}

type rangeFloat64Filter struct {
	Key      string
	From, To float64
}

func (k *rangeFloat64Filter) Check(key, val string) bool {
	if key != k.Key {
		return false
	}
	var (
		floatVal float64
		err      error
	)
	if floatVal, err = strconv.ParseFloat(val, 64); err != nil {
		return false
	}
	return floatVal > k.From && floatVal <= k.To
}
