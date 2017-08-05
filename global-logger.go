package kiwi

// This file consists of definition of global logging methods.

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

// Log is simplified realization of Logger.Log().
// It has no context. You would like use it in short applications where context and
// initialization of logger will look too complex. Use Logger type instead.
func Log(keyVals ...interface{}) {
	var (
		key    string
		record = make([]*Pair, 0, len(keyVals)/2+1)
	)
	for i, val := range keyVals {
		if i%2 == 0 {
			key = toKey(val)
			continue
		}
		var p *Pair
		if p = toPair(key, val); p.Eval != nil {
			p.Val = p.Eval.(func() string)()
		}
		record = append(record, p)
	}
	if len(keyVals)%2 == 1 {
		record = append(record, &Pair{key, "", nil, VoidVal})
	}
	collector.WaitFlush.Add(collector.Count)
	// It will be unlocked inside sinkRecord().
	collector.RLock()
	go sinkRecord(record)
}
