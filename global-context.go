package kiwi

// Global context for all logger instances including global logger.

/* Copyright (c) 2016-2024, Alexander I.Grafov <grafov@inet.name>
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
	"sync"
)

var (
	global  sync.RWMutex
	context []*Pair
)

// With adds key-vals to the global logger context. It is safe for
// concurrency.
func With(kv ...interface{}) {
	var (
		key       string
		thisIsKey = true
	)
	global.Lock()
next:
	for _, arg := range kv {
		if thisIsKey {
			switch p := arg.(type) {
			// The odd arg treated as the keys. The key must be a
			// string.
			case string:
				key = p
			// Instead of the key the key-value pair could be
			// passed. Next arg should be a key.
			case *Pair:
				for i, c := range context {
					if c.Key == p.Key {
						context[i] = p
						break next
					}
				}
				context = append(context, p)
				continue
			// Also the slice of key-value pairs could be passed. Next
			// arg should be a key.
			case []*Pair:
				for _, v := range p {
					for i, c := range context {
						if c.Key == v.Key {
							context[i] = v
							break
						}
					}
					context = append(context, v)
				}
				continue
			// The key must be be a string type. The logger generates
			// error as a new key-value pair for the record.
			default:
				context = append(context, toPair(ErrorKey, "wrong type for the key"))
				key = MessageKey
			}
		} else {
			p := toPair(key, arg)
			for i, c := range context {
				if c.Key == key {
					context[i] = p
					thisIsKey = !thisIsKey
					break next
				}
			}
			context = append(context, p)
		}
		thisIsKey = !thisIsKey
	}
	if !thisIsKey && key != MessageKey {
		context = append(context, toPair(MessageKey, key))
	}
	global.Unlock()
}

// Without drops the keys from the context of the global logger. It is safe for
// concurrency.
func Without(keys ...string) {
	global.Lock()
	for _, key := range keys {
		for i, p := range context {
			if p.Key == key {
				copy(context[i:], context[i+1:])
				context[len(context)-1] = nil
				context = context[:len(context)-1]
				break
			}
		}
	}
	global.Unlock()
}
