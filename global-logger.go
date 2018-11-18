package kiwi

// This file consists of definition of global logging methods.

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

// Log is simplified realization of Logger.Log().
// You would like use it in short applications where context and
// initialization of logger could brought extra complexity.
// If you wish separate contexts and achieve better performance
// use Logger type instead.
func Log(keyVals ...interface{}) {
	var (
		record  = make([]*Pair, 0, len(keyVals)/2+1)
		key     interface{}
		nextKey = true
	)
	globalContext.RLock()
	for _, p := range globalContext.m {
		if p.Eval != nil {
			// Evaluate delayed context value here before output.
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, p.Eval, p.Type})
		}
	}
	globalContext.RUnlock()
	for _, val := range keyVals {
		var p *Pair
		if nextKey {
			switch val.(type) {
			case Pair:
				p = val.(*Pair)
				if p.Eval != nil {
					p.Val = p.Eval.(func() string)()
				}
				record = append(record, p)
				continue
			}
			key = val
			nextKey = false
		} else {
			if p = toPair(key.(string), val); p.Eval != nil {
				p.Val = p.Eval.(func() string)()
			}
			record = append(record, p)
			nextKey = true
		}
	}
	// add the value without the key for odd number for key-val pairs
	if !nextKey {
		record = append(record, toPair(UnpairedKey, key))
	}
	collector.WaitFlush.Add(collector.Count)
	// It will be unlocked inside sinkRecord().
	collector.RLock()
	go sinkRecord(record)
}
