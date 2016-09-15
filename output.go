package kiwi

// This file consists of Output related structures and functions.
// Outputs accepts incoming log records from Loggers, check them with filters
// and write to output streams if checks passed.

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
	"io"
	"sync"
)

var outputs []*Output

type (
	// Output used for filtering incoming log records from all logger instances
	// and decides how to filter them. Each output wraps its own io.Writer.
	// Output methods are safe for concurrent usage.
	Output struct {
		sync.RWMutex
		In              chan *[]pair
		writer          io.Writer
		format          Formatter
		paused          bool
		closed          bool
		positiveFilters map[string]filter
		negativeFilters map[string]filter
		hiddenKeys      map[string]bool
	}
)

// UseOutput creates a new output for an arbitrary number of loggers.
// There are any number of outputs may be created for saving incoming log
// records to different places.
func UseOutput(w io.Writer, fn Formatter) *Output {
	for i, output := range outputs {
		if output.writer == w {
			outputs[i].format = fn
			return outputs[i]
		}
	}
	output := &Output{
		In:              make(chan *[]pair, 16),
		writer:          w,
		positiveFilters: make(map[string]filter),
		negativeFilters: make(map[string]filter),
		format:          fn,
	}
	outputs = append(outputs, output)
	go processOutput(output)
	return output
}

// With sets restriction for records output.
// Only the records WITH any of the keys will be passed to output.
func (o *Output) With(keys ...string) *Output {
	if !o.closed {
		o.Lock()
		for _, tag := range keys {
			delete(o.negativeFilters, tag)
			o.positiveFilters[tag] = &keyFilter{Key: tag}
		}
		o.Unlock()
	}
	return o
}

// Without sets restriction for records output.
// Only the records WITHOUT any of the keys will be passed to output.
func (o *Output) Without(keys ...string) *Output {
	if !o.closed {
		o.Lock()
		for _, tag := range keys {
			o.negativeFilters[tag] = &keyFilter{Key: tag}
			delete(o.positiveFilters, tag)
		}
		o.Unlock()
	}
	return o
}

// WithValues sets restriction for records output.
// A record passed to output if the key equal one of any of the listed values.
func (o *Output) WithValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return o.With(key)
	}
	if !o.closed {
		o.Lock()
		delete(o.negativeFilters, key)
		o.positiveFilters[key] = &valsFilter{Key: key, Vals: vals}
		o.Unlock()
	}
	return o
}

// WithoutValues sets restriction for records output.
func (o *Output) WithoutValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return o.Without(key)
	}
	if !o.closed {
		o.Lock()
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &valsFilter{Key: key, Vals: vals}
		o.Unlock()
	}
	return o
}

// WithRangeInt64 sets restriction for records output.
func (o *Output) WithRangeInt64(key string, from, to int64) *Output {
	if !o.closed {
		o.Lock()
		delete(o.negativeFilters, key)
		o.positiveFilters[key] = &rangeInt64Filter{Key: key, From: from, To: to}
		o.Unlock()
	}
	return o
}

// WithRangeFloat64 sets restriction for records output.
func (o *Output) WithoutRangeInt64(key string, from, to int64) *Output {
	o.Lock()
	if !o.closed {
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &rangeInt64Filter{Key: key, From: from, To: to}
	}
	o.Unlock()
	return o
}

// WithRangeFloat64 sets restriction for records output.
func (o *Output) WithRangeFloat64(key string, from, to float64) *Output {
	o.Lock()
	delete(o.negativeFilters, key)
	o.positiveFilters[key] = &rangeFloat64Filter{Key: key, From: from, To: to}
	o.Unlock()
	return o
}

// WithoutRangeFloat64  sets restriction for records output.
func (o *Output) WithoutRangeFloat64(key string, from, to float64) *Output {
	if !o.closed {
		o.Lock()
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &rangeFloat64Filter{Key: key, From: from, To: to}
		o.Unlock()
	}
	return o
}

// Reset all filters for the keys for the output.
func (o *Output) Reset(keys ...string) *Output {
	o.Lock()
	for _, tag := range keys {
		delete(o.positiveFilters, tag)
		delete(o.negativeFilters, tag)
	}
	o.Unlock()
	return o
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (o *Output) Hide(keys ...string) *Output {
	o.Lock()
	if !o.closed {
		for _, tag := range keys {
			o.hiddenKeys[tag] = true
		}
	}
	o.Unlock()
	return o
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (o *Output) Unhide(keys ...string) *Output {
	o.Lock()
	if !o.closed {
		for _, tag := range keys {
			delete(o.hiddenKeys, tag)
		}
	}
	o.Unlock()
	return o
}

// Pause stops writing to the output.
func (o *Output) Pause() {
	o.paused = true
}

// Contiunue writing to the output.
func (o *Output) Continue() {
	o.paused = false
}

func (o *Output) Close() {
	o.Lock()
	o.closed = true
	o.Unlock()
	close(o.In)
}

// Flush waits that all previously sent to the output records worked.
func (o *Output) Flush() {
	var flush = make(chan struct{})
	// Well, it uses some kind of lifehack instead of dedicated flag.
	// It send "deleted" record with unbuffered channel in the value.
	// Then just wait for the value from this channel.
	o.In <- &[]pair{{Deleted: true, Val: value{Func: flush}}}
	<-flush
}

func processOutput(o *Output) {
	for {
		record, ok := <-o.In
		if !ok {
			o.positiveFilters = nil
			o.negativeFilters = nil
			o.hiddenKeys = nil
			return
		}
		if o.closed || o.paused {
			continue
		}
		o.RLock()
		for i, pair := range *record {
			// Flush came
			if i == 0 && pair.Deleted {
				pair.Val.Func.(chan struct{}) <- struct{}{}
				goto skipRecord
			}
			// Negative conditions have highest priority
			if filter, ok := o.negativeFilters[pair.Key]; ok {
				if filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
			// At last check for positive conditions
			if filter, ok := o.positiveFilters[pair.Key]; ok {
				if !filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
		}
		o.filter(record)
		continue
	skipRecord:
		o.RUnlock()
	}
}

func (o *Output) filter(record *[]pair) {
	o.format.Begin()
	for _, pair := range *record {
		if ok := o.hiddenKeys[pair.Key]; ok {
			continue
		}
		o.format.Pair(pair.Key, pair.Val.Strv, pair.Val.Quoted)
	}
	o.RUnlock()
	o.writer.Write(o.format.Finish())
}
