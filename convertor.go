// Convert incoming values to string representation. For keys and values.
package kiwi

import (
	"fmt"
	"strconv"
)

const (
	voidVal uint8 = iota
	stringVal
	booleanVal
	integerVal
	floatVal
	complexVal
)

// it applicable for all scalar types and for strings
func toRecordKey(val interface{}) string {
	switch val.(type) {
	case bool:
		if val.(bool) {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(val.(int))
	case int8:
		return strconv.FormatInt(int64(val.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(val.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	case uint:
		return strconv.FormatUint(uint64(val.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case float32:
		return strconv.FormatFloat(float64(val.(float32)), 'e', -1, 32)
	case float64:
		return strconv.FormatFloat(val.(float64), 'e', -1, 64)
	case complex64:
		return fmt.Sprintf("%f", val.(complex64))
	case complex128:
		return fmt.Sprintf("%f", val.(complex128))
	case fmt.Stringer:
		return val.(fmt.Stringer).String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// it applicable for all scalar types and for strings
func toRecordValue(val interface{}) value {
	switch val.(type) {
	case string:
		return value{val.(string), nil, stringVal, true}
	case bool:
		if val.(bool) {
			return value{"true", nil, booleanVal, false}
		}
		return value{"false", nil, booleanVal, false}
	case int:
		return value{strconv.Itoa(val.(int)), nil, integerVal, false}
	case int8:
		return value{strconv.FormatInt(int64(val.(int8)), 10), nil, integerVal, false}
	case int16:
		return value{strconv.FormatInt(int64(val.(int16)), 10), nil, integerVal, false}
	case int32:
		return value{strconv.FormatInt(int64(val.(int32)), 10), nil, integerVal, false}
	case int64:
		return value{strconv.FormatInt(val.(int64), 10), nil, integerVal, false}
	case uint:
		return value{strconv.FormatUint(uint64(val.(uint)), 10), nil, integerVal, false}
	case uint8:
		return value{strconv.FormatUint(uint64(val.(uint8)), 10), nil, integerVal, false}
	case uint16:
		return value{strconv.FormatUint(uint64(val.(uint16)), 10), nil, integerVal, false}
	case uint32:
		return value{strconv.FormatUint(uint64(val.(uint32)), 10), nil, integerVal, false}
	case uint64:
		return value{strconv.FormatUint(val.(uint64), 10), nil, integerVal, false}
	case float32:
		return value{strconv.FormatFloat(float64(val.(float32)), 'e', -1, 32), nil, floatVal, false}
	case float64:
		return value{strconv.FormatFloat(val.(float64), 'e', -1, 64), nil, floatVal, false}
	case complex64:
		return value{fmt.Sprintf("%f", val.(complex64)), nil, complexVal, false}
	case complex128:
		return value{fmt.Sprintf("%f", val.(complex128)), nil, complexVal, false}
	case Record:
		return value{val.(Record).String(), nil, stringVal, true}
	case fmt.Stringer:
		return value{val.(fmt.Stringer).String(), nil, stringVal, true}
	case func() string:
		return value{"", val, stringVal, true}
	case func() bool:
		return value{"", val, booleanVal, true}
	case func() int, func() int8, func() int16, func() int32, func() int64:
		return value{"", val, integerVal, true}
	case func() uint8, func() uint16, func() uint32, func() uint64:
		return value{"", val, integerVal, true}
	case func() float32, func() float64:
		return value{"", val, floatVal, true}
	case func() complex64:
		return value{"", val, complexVal, false}
	case nil:
		return value{"", nil, voidVal, false}
	default:
		return value{fmt.Sprintf("%v", val), nil, stringVal, true}
	}
}

// calls function()T return its result as an interface
func toFunc(fn interface{}) interface{} {
	switch fn.(type) {
	case func() string:
		return fn.(func() string)()
	}
	return nil
}
