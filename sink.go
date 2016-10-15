package kiwi

// This file consists of Sink related structures and functions.
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
	"time"
)

// Sinks accepts records through the chanels.
// Each sink has its own channel.
var sinks []*Sink

type (
	// Sink used for filtering incoming log records from all logger instances
	// and decides how to filter them. Each output wraps its own io.Writer.
	// Sink methods are safe for concurrent usage.
	Sink struct {
		sync.RWMutex
		In              chan *[]pair
		writer          io.Writer
		format          Formatter
		paused          bool
		closed          bool
		positiveFilters map[string]Filter
		negativeFilters map[string]Filter
		hiddenKeys      map[string]bool
	}
)

// SinkTo creates a new sink for an arbitrary number of loggers.
// There are any number of sinks may be created for saving incoming log
// records to different places.
func SinkTo(w io.Writer, fn Formatter) *Sink {
	for i, output := range sinks {
		if output.writer == w {
			sinks[i].format = fn
			return sinks[i]
		}
	}
	output := &Sink{
		In:              make(chan *[]pair, 16),
		writer:          w,
		positiveFilters: make(map[string]Filter),
		negativeFilters: make(map[string]Filter),
		format:          fn,
		paused:          true, // it started paused because not pass records until the filters set
	}
	sinks = append(sinks, output)
	go processOutput(output)
	return output
}

// WithKey sets restriction for records output.
// Only the records WITH any of the keys will be passed to output.
func (o *Sink) WithKey(keys ...string) *Sink {
	if !o.closed {
		o.Lock()
		for _, key := range keys {
			o.positiveFilters[key] = &keyFilter{}
			delete(o.negativeFilters, key)
		}
		o.Unlock()
	}
	return o
}

// WithoutKey sets restriction for records output.
// Only the records WITHOUT any of the keys will be passed to output.
func (o *Sink) WithoutKey(keys ...string) *Sink {
	if !o.closed {
		o.Lock()
		for _, key := range keys {
			o.negativeFilters[key] = &keyFilter{}
			delete(o.positiveFilters, key)
		}
		o.Unlock()
	}
	return o
}

// WithValue sets restriction for records output.
// A record passed to output if the key equal one of any of the listed values.
func (o *Sink) WithValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return o.WithKey(key)
	}
	if !o.closed {
		o.Lock()
		o.positiveFilters[key] = &valsFilter{Vals: vals}
		delete(o.negativeFilters, key)
		o.Unlock()
	}
	return o
}

// WithoutValue sets restriction for records output.
func (o *Sink) WithoutValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return o.WithoutKey(key)
	}
	if !o.closed {
		o.Lock()
		o.negativeFilters[key] = &valsFilter{Vals: vals}
		delete(o.positiveFilters, key)
		o.Unlock()
	}
	return o
}

// WithInt64Range sets restriction for records output.
func (o *Sink) WithInt64Range(key string, from, to int64) *Sink {
	if !o.closed {
		o.Lock()
		delete(o.negativeFilters, key)
		o.positiveFilters[key] = &int64RangeFilter{From: from, To: to}
		o.Unlock()
	}
	return o
}

// WithoutInt64Range sets restriction for records output.
func (o *Sink) WithoutInt64Range(key string, from, to int64) *Sink {
	o.Lock()
	if !o.closed {
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &int64RangeFilter{From: from, To: to}
	}
	o.Unlock()
	return o
}

// WithFloat64Range sets restriction for records output.
func (o *Sink) WithFloat64Range(key string, from, to float64) *Sink {
	o.Lock()
	delete(o.negativeFilters, key)
	o.positiveFilters[key] = &float64RangeFilter{From: from, To: to}
	o.Unlock()
	return o
}

// WithoutFloat64Range sets restriction for records output.
func (o *Sink) WithoutFloat64Range(key string, from, to float64) *Sink {
	if !o.closed {
		o.Lock()
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &float64RangeFilter{From: from, To: to}
		o.Unlock()
	}
	return o
}

// WithTimeRange sets restriction for records output.
func (o *Sink) WithTimeRange(key string, from, to time.Time) *Sink {
	o.Lock()
	delete(o.negativeFilters, key)
	o.positiveFilters[key] = &timeRangeFilter{From: from, To: to}
	o.Unlock()
	return o
}

// WithoutTimeRange sets restriction for records output.
func (o *Sink) WithoutTimeRange(key string, from, to time.Time) *Sink {
	if !o.closed {
		o.Lock()
		delete(o.positiveFilters, key)
		o.negativeFilters[key] = &timeRangeFilter{From: from, To: to}
		o.Unlock()
	}
	return o
}

// WithFilter setup custom filtering function for values for a specific key.
// Custom filter should realize Filter interface. All custom filters treated
// as positive filters. So if the filter returns true then it will be passed.
func (o *Sink) WithFilter(key string, customFilter Filter) *Sink {
	if !o.closed {
		o.Lock()
		o.positiveFilters[key] = customFilter
		o.Unlock()
	}
	return o
}

// Reset all filters for the keys for the output.
func (o *Sink) Reset(keys ...string) *Sink {
	o.Lock()
	for _, key := range keys {
		delete(o.positiveFilters, key)
		delete(o.negativeFilters, key)
	}
	o.Unlock()
	return o
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (o *Sink) Hide(keys ...string) *Sink {
	o.Lock()
	if !o.closed {
		for _, key := range keys {
			o.hiddenKeys[key] = true
		}
	}
	o.Unlock()
	return o
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (o *Sink) Unhide(keys ...string) *Sink {
	o.Lock()
	if !o.closed {
		for _, key := range keys {
			delete(o.hiddenKeys, key)
		}
	}
	o.Unlock()
	return o
}

// Stop stops writing to the output.
func (o *Sink) Stop() *Sink {
	o.paused = true
	return o
}

// Start writing to the output.
// After creation of a new sink it will paused and you need explicitly start it.
func (o *Sink) Start() *Sink {
	o.paused = false
	return o
}

// Close the sink. Flush all records that came before.
func (o *Sink) Close() {
	if o.closed {
		return
	}
	o.Lock()
	o.closed = true
	o.Unlock()
	close(o.In)
}

// Flush waits that all previously sent to the output records worked.
func (o *Sink) Flush() *Sink {
	if !o.paused && !o.closed {
		var flush = make(chan struct{})
		// Well, it uses some kind of lifehack instead of dedicated flag.
		// It send "deleted" record with unbuffered channel in the value.
		// Then just wait for the value from this channel.
		o.In <- &[]pair{{Deleted: true, Val: value{Func: flush}}}
		<-flush
	}
	return o
}

func processOutput(o *Sink) {
	for {
		record, ok := <-o.In
		if !ok {
			o.positiveFilters = nil
			o.negativeFilters = nil
			o.hiddenKeys = nil
			o.closed = true
			return
		}
		if o.paused || o.closed {
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

func (o *Sink) filter(record *[]pair) {
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
