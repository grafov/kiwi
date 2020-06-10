package kiwi

// Convert incoming values to string representation. For keys and values.

/* Copyright (c) 2016-2020, Alexander I.Grafov <grafov@gmail.com>
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

// Possible kinds of logged values.
const (
	// BooleanVal and other types below commonly formatted unquoted.
	// But it depends on the formatter.
	BooleanVal = iota
	IntegerVal
	FloatVal
	ComplexVal
	CustomUnquoted
	// VoidVal and other types below commonly formatted unquoted.
	// But it depends on the formatter.
	StringVal
	TimeVal
	CustomQuoted
)

// FloatFormat used in Float to String conversion.
// It is second parameter passed to strconv.FormatFloat()
var FloatFormat byte = 'e'

// TimeLayout used in time.Time to String conversion.
var TimeLayout = time.RFC3339

// it applicable for all scalar types and for strings
func toPair(key string, val interface{}) *Pair {
	switch v := val.(type) {
	case string:
		return &Pair{key, v, nil, StringVal}
	case []byte:
		return &Pair{key, string(v), nil, StringVal}
	case bool:
		if val.(bool) {
			return &Pair{key, "true", nil, BooleanVal}
		}
		return &Pair{key, "false", nil, BooleanVal}
	case int:
		return &Pair{key, strconv.Itoa(v), nil, IntegerVal}
	case int8:
		return &Pair{key, strconv.FormatInt(int64(v), 10), nil, IntegerVal}
	case int16:
		return &Pair{key, strconv.FormatInt(int64(v), 10), nil, IntegerVal}
	case int32:
		return &Pair{key, strconv.FormatInt(int64(v), 10), nil, IntegerVal}
	case int64:
		return &Pair{key, strconv.FormatInt(v, 10), nil, IntegerVal}
	case uint:
		return &Pair{key, strconv.FormatUint(uint64(v), 10), nil, IntegerVal}
	case uint8:
		return &Pair{key, strconv.FormatUint(uint64(v), 10), nil, IntegerVal}
	case uint16:
		return &Pair{key, strconv.FormatUint(uint64(v), 10), nil, IntegerVal}
	case uint32:
		return &Pair{key, strconv.FormatUint(uint64(v), 10), nil, IntegerVal}
	case uint64:
		return &Pair{key, strconv.FormatUint(v, 10), nil, IntegerVal}
	case float32:
		return &Pair{key, strconv.FormatFloat(float64(v), FloatFormat, -1, 32), nil, FloatVal}
	case float64:
		return &Pair{key, strconv.FormatFloat(v, FloatFormat, -1, 64), nil, FloatVal}
	case complex64:
		return &Pair{key, fmt.Sprintf("%f", v), nil, ComplexVal}
	case complex128:
		return &Pair{key, fmt.Sprintf("%f", v), nil, ComplexVal}
	case time.Time:
		return &Pair{key, v.Format(TimeLayout), nil, TimeVal}
	case Valuer:
		var pairType = CustomUnquoted
		if v.IsQuoted() {
			pairType = CustomQuoted
		}
		return &Pair{key, v.String(), nil, pairType}
	case Stringer:
		return &Pair{key, v.String(), nil, StringVal}
	case encoding.TextMarshaler:
		data, err := v.MarshalText()
		if err != nil {
			return &Pair{key, fmt.Sprintf("%s", err), nil, StringVal}
		}
		return &Pair{key, string(data), nil, StringVal}
	case func() string:
		return &Pair{key, "", v, StringVal}
	default:
		// Worst case conversion that depends on reflection.
		return &Pair{key, fmt.Sprintf("%+v", val), nil, StringVal}
	}
}
