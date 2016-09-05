package kiwi

/*
Copyright (c) 2016, Alexander I.Grafov aka Axel
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

ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"fmt"
	"strconv"
)

func String(key string, val string) pair {
	return pair{key, value{val, nil, stringVal, true}, false}
}

func Stringer(key string, val fmt.Stringer) pair {
	return pair{key, value{val.String(), nil, stringVal, true}, false}
}

func Int(key string, val int) pair {
	return pair{key, value{strconv.Itoa(val), nil, integerVal, true}, false}
}

func Int64(key string, val int64) pair {
	return pair{key, value{strconv.FormatInt(val, 10), nil, integerVal, true}, false}
}

func Float64(key string, val float64) pair {
	return pair{key, value{strconv.FormatFloat(val, 'e', -1, 64), nil, floatVal, true}, false}
}

func Bool(key string, val bool) pair {
	if val {
		return pair{key, value{"true", nil, booleanVal, false}, false}
	}
	return pair{key, value{"false", nil, booleanVal, false}, false}
}
