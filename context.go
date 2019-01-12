package kiwi

// The logger instance context.

/* Copyright (c) 2016-2019, Alexander I.Grafov <grafov@gmail.com>
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

import "fmt"

// With defines a context for the logger. The context overrides pairs
// in the record. The function is not concurrent safe.
func (l *Logger) With(keyVals ...interface{}) *Logger {
	var (
		key          string
		shouldBeAKey = true
	)
	// key=val pairs
	for _, val := range keyVals {
		if shouldBeAKey {
			switch val.(type) {
			case string:
				key = val.(string)
			case *Pair:
				p := val.(*Pair)
				l.context[p.Key] = p
				continue
			case []*Pair:
				for _, p := range val.([]*Pair) {
					l.context[p.Key] = p
				}
				continue
			default:
				l.context[ErrorKey] = toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val))
				key = UnpairedKey
			}
		} else {
			l.context[key] = toPair(key, val)
		}
		shouldBeAKey = !shouldBeAKey
	}
	if !shouldBeAKey && key != UnpairedKey {
		l.pairs = append(l.pairs, toPair(UnpairedKey, key))
	}
	return l
}

// Without drops some keys from a context for the logger. The function
// is not concurrent safe.
func (l *Logger) Without(keys ...string) *Logger {
	for _, key := range keys {
		delete(l.context, key)
	}
	return l
}

// ResetContext resets the context of the logger. The function is not
// concurrent safe.
func (l *Logger) ResetContext() *Logger {
	l.context = make(map[string]*Pair, len(l.context)*2)
	return l
}
