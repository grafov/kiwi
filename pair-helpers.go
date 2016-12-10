package kiwi

// This file consists of helpers for adding pairs with strongly typed values.
// I not sure about this part of API yet. It should be moved into a separate
// package or removed.

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
	"strconv"
	"time"
)

// AsString formats pair for string.
// Note: type helpers are experimental part of API and may be removed.
func AsString(key string, val string) *pair {
	return &pair{key, val, nil, StringVal}
}

// AsStringer formats pair for string with Stringer interface (the same as fmt.Stringer).
// Note: type helpers are experimental part of API and may be removed.
func AsStringer(key string, val Stringer) *pair {
	return &pair{key, val.String(), nil, StringVal}
}

// AsInt formats pair for int value. If you need add integer of specific size just
// convert it to int, int64 or uint64 and use AddInt(), AddInt64() or AddUint64()
// respectively.
// Note: type helpers are experimental part of API and may be removed.
func AsInt(key string, val int) *pair {
	return &pair{key, strconv.Itoa(val), nil, IntegerVal}
}

// AsInt64 formats pair for int64 value.
// Note: type helpers are experimental part of API and may be removed.
func AsInt64(key string, val int64) *pair {
	return &pair{key, strconv.FormatInt(val, 10), nil, IntegerVal}
}

// AsUint64 formats pair for uint64 value.
// Note: type helpers are experimental part of API and may be removed.
func AsUint64(key string, val uint64) *pair {
	return &pair{key, strconv.FormatUint(val, 10), nil, IntegerVal}
}

// AsFloat64 formats pair for float64 value. If you need add float of other size just
// convert it to float64.
// Note: type helpers are experimental part of API and may be removed.
func AsFloat64(key string, val float64) *pair {
	return &pair{key, strconv.FormatFloat(val, 'e', -1, 64), nil, FloatVal}
}

// AsBool formats pair for bool value.
// Note: type helpers are experimental part of API and may be removed.
func AsBool(key string, val bool) *pair {
	if val {
		return &pair{key, "true", nil, BooleanVal}
	}
	return &pair{key, "false", nil, BooleanVal}
}

// AsTime formats pair for time.Time value.
// Note: type helpers are experimental part of API and may be removed.
func AsTime(key string, val time.Time, layout string) *pair {
	return &pair{key, val.Format(layout), nil, TimeVal}
}
