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
		key       string
		thisIsKey = true
	)
next:
	for _, arg := range keyVals {
		if thisIsKey {
			switch arg.(type) {
			// The odd arg treated as the keys. The key must be a
			// string.
			case string:
				key = arg.(string)
			// Instead of the key the key-value pair could be
			// passed. Next arg should be a key.
			case *Pair:
				p := arg.(*Pair)
				for i, c := range l.context {
					if p.Key == c.Key {
						l.context[i] = p
						break next
					}
				}
				l.context = append(l.context, p)
				continue
			// Also the slice of key-value pairs could be passed. Next
			// arg should be a key.
			case []*Pair:
				for _, p := range arg.([]*Pair) {
					for i, c := range l.context {
						if c.Key == p.Key {
							l.context[i] = p
							break
						}
					}
					l.context = append(l.context, p)
				}
				continue
			// The key must be be a string type. The logger generates
			// error as a new key-value pair for the record.
			default:
				l.context = append(l.context, toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", arg, arg)))
				key = UnpairedKey
			}
		} else {
			p := toPair(key, arg)
			for i, c := range l.context {
				if c.Key == key {
					l.context[i] = p
					thisIsKey = !thisIsKey
					break next
				}
			}
			l.context = append(l.context, toPair(key, arg))
		}
		// After the key the next arg is not a key.
		thisIsKey = !thisIsKey
	}
	if !thisIsKey && key != UnpairedKey {
		l.context = append(l.context, toPair(UnpairedKey, key))
	}
	return l
}

// Without drops some keys from a context for the logger. The function
// is not concurrent safe.
func (l *Logger) Without(keys ...string) *Logger {
	for _, k := range keys {
		for i, v := range l.context {
			if v.Key == k {
				copy(l.context[i:], l.context[i+1:])
				l.context[len(l.context)-1] = nil
				l.context = l.context[:len(l.context)-1]
				break
			}
		}
	}
	return l
}

// ResetContext resets the context of the logger. The function is not
// concurrent safe.
func (l *Logger) ResetContext() *Logger {
	l.context = nil
	return l
}
