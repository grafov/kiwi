package kiwi

import (
	"sync"
)

// Global context for all logger instances including global logger.

/* Copyright (c) 2016-2018, Alexander I.Grafov <grafov@gmail.com>
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

type context struct {
	sync.RWMutex
	m map[string]Pair
}

var globalContext context

func With(keyVals ...interface{}) {
	var (
		key     interface{}
		keyStr  string
		nextKey = true
	)
	// key=val pairs
	globalContext.Lock()
	for _, val := range keyVals {
		if nextKey {
			switch val.(type) {
			case Pair:
				globalContext.m[keyStr] = val.(Pair)
				continue
			case []*Pair:
				for _, p := range val.([]*Pair) {
					globalContext.m[p.Key] = *p
				}
				continue
			}
			key = val
			nextKey = false
		} else {
			keyStr = key.(string)
			globalContext.m[keyStr] = *toPair(keyStr, val)
			nextKey = true
		}
	}
	//  add the value without the key for odd number for key-val pairs
	if !nextKey {
		globalContext.m[UnpairedKey] = *toPair(UnpairedKey, key)
	}
	globalContext.Unlock()
}

// Without drops some keys from a context for the logger.
func Without(keys ...string) {
	globalContext.Lock()
	for _, key := range keys {
		delete(globalContext.m, key)
	}
	globalContext.Unlock()
}

// ResetContext resets the global context for the global logger and
// its descendants.
func ResetContext() {
	globalContext.Lock()
	globalContext.m = make(map[string]Pair, len(globalContext.m))
	globalContext.Unlock()
}

func init() {
	globalContext.m = make(map[string]Pair)
}
