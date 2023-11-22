package kiwi

import "fmt"

// This file consists of definition of global logging methods.

/* Copyright (c) 2016-2019, 2023, Alexander I.Grafov <grafov@inet.name>
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

// Log is simplified realization of Logger.Log().
// You would like use it in short applications where context and
// initialization of logger could brought extra complexity.
// If you wish separate contexts and achieve better performance
// use Logger type instead.
func Log(kv ...interface{}) {
	// 1. Log the context.
	record := make([]*Pair, 0, len(context)+len(kv))
	global.RLock()
	for _, p := range context {
		// Evaluate delayed context value here before the output.
		if p.Eval != nil {
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, nil, p.Type})
		}
	}
	global.RUnlock()
	// 2. Log the regular key-value pairs that came in the args.
	var (
		key         string
		shouldBeKey = true
	)
	for _, val := range kv {
		var p *Pair
		if shouldBeKey {
			switch v := val.(type) {
			case string:
				key = v
			case *Pair:
				if v.Eval != nil {
					v.Val = v.Eval.(func() string)()
				}
				record = append(record, v)
				continue
			default:
				record = append(record, toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val)))
				key = MessageKey
			}
		} else {
			if p = toPair(key, val); p.Eval != nil {
				p.Val = p.Eval.(func() string)()
			}
			record = append(record, p)
		}
		shouldBeKey = !shouldBeKey
	}
	// Add the value without the key for odd number for key-val pairs.
	if !shouldBeKey && key != MessageKey {
		record = append(record, toPair(MessageKey, key))
	}
	// 2. Pass the record to the collector.
	sinkRecord(record)
}
