package strict

// This file consists of helpers for adding pairs with strongly typed values.
// I not sure about this part of API yet. It should be moved into a separate
// package or removed.

/* Copyright (c) 2016-2017, Alexander I.Grafov
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

	"github.com/grafov/kiwi"
)

// String formats pair for string.
// Note: type helpers are experimental part of API and may be removed.
func String(key string, val string) *kiwi.Pair {
	return &kiwi.Pair{key, val, nil, kiwi.StringVal}
}

// Stringer formats pair for string.
// Note: type helpers are experimental part of API and may be removed.
func Stringer(key string, val kiwi.Stringer) *kiwi.Pair {
	return &kiwi.Pair{key, val.String(), nil, kiwi.StringVal}
}

// Int formats pair for int value. If you need add integer of specific size just
// convert it to int, int64 or uint64 and use AddInt(), AddInt64() or AddUint64()
// respectively.
// Note: type helpers are experimental part of API and may be removed.
func Int(key string, val int) *kiwi.Pair {
	return &kiwi.Pair{key, strconv.Itoa(val), nil, kiwi.IntegerVal}
}

// Int64 formats pair for int64 value.
// Note: type helpers are experimental part of API and may be removed.
func Int64(key string, val int64) *kiwi.Pair {
	return &kiwi.Pair{key, strconv.FormatInt(val, 10), nil, kiwi.IntegerVal}
}

// Uint64 formats pair for uint64 value.
// Note: type helpers are experimental part of API and may be removed.
func Uint64(key string, val uint64) *kiwi.Pair {
	return &kiwi.Pair{key, strconv.FormatUint(val, 10), nil, kiwi.IntegerVal}
}

// Float64 formats pair for float64 value. If you need add float of other size just
// convert it to float64.
// Note: type helpers are experimental part of API and may be removed.
func Float64(key string, val float64) *kiwi.Pair {
	return &kiwi.Pair{key, strconv.FormatFloat(val, 'e', -1, 64), nil, kiwi.FloatVal}
}

// Bool formats pair for bool value.
// Note: type helpers are experimental part of API and may be removed.
func Bool(key string, val bool) *kiwi.Pair {
	if val {
		return &kiwi.Pair{key, "true", nil, kiwi.BooleanVal}
	}
	return &kiwi.Pair{key, "false", nil, kiwi.BooleanVal}
}

// Time formats pair for time.Time value.
// Note: type helpers are experimental part of API and may be removed.
func Time(key string, val time.Time, layout string) *kiwi.Pair {
	return &kiwi.Pair{key, val.Format(layout), nil, kiwi.TimeVal}
}
