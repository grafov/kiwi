package kiwi

// Convert incoming values to string representation. For keys and values.

/* Copyright (c) 2016, Alexander I.Grafov aka Axel
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of kvlog nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

ॐ तारे तुत्तारे तुरे स्व */

import (
	"encoding"
	"fmt"
	"strconv"
	"time"
)

const (
	voidVal uint8 = iota
	stringVal
	booleanVal
	integerVal
	floatVal
	complexVal
	flushCmd
)

// FloatFormat used in Float to String conversion.
// It is second parameter passed to strconv.FormatFloat()
var FloatFormat byte = 'e'

// TimeLayout used in time.Time to String conversion.
var TimeLayout = time.RFC3339

// it applicable for all scalar types and for strings
func toKey(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case []byte:
		return string(val.([]byte))
	case fmt.Stringer:
		return val.(fmt.Stringer).String()
	case encoding.TextMarshaler:
		data, err := val.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		return string(data)
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
		return strconv.FormatFloat(float64(val.(float32)), FloatFormat, -1, 32)
	case float64:
		return strconv.FormatFloat(val.(float64), FloatFormat, -1, 64)
	case complex64:
		return fmt.Sprintf("%f", val.(complex64))
	case complex128:
		return fmt.Sprintf("%f", val.(complex128))
	case time.Time:
		return val.(time.Time).Format(TimeLayout)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// it applicable for all scalar types and for strings
func toPair(key string, val interface{}) *pair {
	switch val.(type) {
	case string:
		return newPair(key, val.(string), nil, stringVal, true)
	case []byte:
		return newPair(key, string(val.([]byte)), nil, stringVal, true)
	case bool:
		if val.(bool) {
			return newPair(key, "true", nil, booleanVal, false)
		}
		return newPair(key, "false", nil, booleanVal, false)
	case int:
		return newPair(key, strconv.Itoa(val.(int)), nil, integerVal, false)
	case int8:
		return newPair(key, strconv.FormatInt(int64(val.(int8)), 10), nil, integerVal, false)
	case int16:
		return newPair(key, strconv.FormatInt(int64(val.(int16)), 10), nil, integerVal, false)
	case int32:
		return newPair(key, strconv.FormatInt(int64(val.(int32)), 10), nil, integerVal, false)
	case int64:
		return newPair(key, strconv.FormatInt(val.(int64), 10), nil, integerVal, false)
	case uint:
		return newPair(key, strconv.FormatUint(uint64(val.(uint)), 10), nil, integerVal, false)
	case uint8:
		return newPair(key, strconv.FormatUint(uint64(val.(uint8)), 10), nil, integerVal, false)
	case uint16:
		return newPair(key, strconv.FormatUint(uint64(val.(uint16)), 10), nil, integerVal, false)
	case uint32:
		return newPair(key, strconv.FormatUint(uint64(val.(uint32)), 10), nil, integerVal, false)
	case uint64:
		return newPair(key, strconv.FormatUint(val.(uint64), 10), nil, integerVal, false)
	case float32:
		return newPair(key, strconv.FormatFloat(float64(val.(float32)), FloatFormat, -1, 32), nil, floatVal, false)
	case float64:
		return newPair(key, strconv.FormatFloat(val.(float64), FloatFormat, -1, 64), nil, floatVal, false)
	case complex64:
		return newPair(key, fmt.Sprintf("%f", val.(complex64)), nil, complexVal, false)
	case complex128:
		return newPair(key, fmt.Sprintf("%f", val.(complex128)), nil, complexVal, false)
	case time.Time:
		return newPair(key, val.(time.Time).Format(TimeLayout), nil, stringVal, false)
	case Valuer:
		return newPair(key, val.(Valuer).String(), nil, stringVal, true)
	case Stringer:
		return newPair(key, val.(Stringer).String(), nil, stringVal, true)
	case encoding.TextMarshaler:
		data, err := val.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return newPair(key, fmt.Sprintf("%s", err), nil, stringVal, true)
		}
		return newPair(key, string(data), nil, stringVal, true)
	case func() string:
		return newPair(key, "", val.(func() string), stringVal, true)
	case nil:
		return newPair(key, "", nil, voidVal, false)
	default:
		// Worst case conversion that depends on reflection.
		return newPair(key, fmt.Sprintf("%v", val), nil, stringVal, true)
	}
}
