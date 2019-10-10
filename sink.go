package kiwi

// This file consists of Sink related structures and functions.
// Outputs accepts incoming log records from Loggers, check them with filters
// and write to output streams if checks passed.

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

import (
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// States of the sink.
const (
	sinkClosed int32 = iota - 1
	sinkStopped
	sinkActive
)

// Sinks accepts records through the chanels.
// Each sink has its own channel.
var collector struct {
	sync.RWMutex
	sinks []*Sink
	count int
}

type (
	// Sink used for filtering incoming log records from all logger instances
	// and decides how to filter them. Each output wraps its own io.Writer.
	// Sink methods are safe for concurrent usage.
	Sink struct {
		id     uint
		In     chan chain
		close  chan struct{}
		writer io.Writer
		format Formatter
		state  *int32

		sync.RWMutex
		positiveFilters map[string]Filter
		negativeFilters map[string]Filter
		hiddenKeys      map[string]bool
	}
	chain struct {
		wg    *sync.WaitGroup
		pairs []*Pair
	}
)

// SinkTo creates a new sink for an arbitrary number of loggers.
// There are any number of sinks may be created for saving incoming log
// records to different places.
// The sink requires explicit start with Start() before usage.
// That allows firstly setup filters before sink will really accept any records.
func SinkTo(w io.Writer, fn Formatter) *Sink {
	collector.RLock()
	for i, sink := range collector.sinks {
		if sink.writer == w {
			collector.sinks[i].format = fn
			collector.RUnlock()
			return collector.sinks[i]
		}
	}
	collector.RUnlock()
	var (
		state = sinkStopped
		sink  = &Sink{
			In:              make(chan chain, 16),
			close:           make(chan struct{}),
			format:          fn,
			state:           &state,
			writer:          w,
			positiveFilters: make(map[string]Filter),
			negativeFilters: make(map[string]Filter),
			hiddenKeys:      make(map[string]bool),
		}
	)
	collector.Lock()
	sink.id = uint(collector.count)
	collector.sinks = append(collector.sinks, sink)
	collector.count++
	collector.Unlock()
	go processSink(sink)
	return sink
}

// HasKey sets restriction for records output.
// Only the records WITH any of the keys will be passed to output.
func (s *Sink) HasKey(keys ...string) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		for _, key := range keys {
			s.positiveFilters[key] = &keyFilter{}
			delete(s.negativeFilters, key)
		}
		s.Unlock()
	}
	return s
}

// HasNotKey sets restriction for records output.
// Only the records WITHOUT any of the keys will be passed to output.
func (s *Sink) HasNotKey(keys ...string) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		for _, key := range keys {
			s.negativeFilters[key] = &keyFilter{}
			delete(s.positiveFilters, key)
		}
		s.Unlock()
	}
	return s
}

// HasValue sets restriction for records output.
// A record passed to output if the key equal one of any of the listed values.
func (s *Sink) HasValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return s.HasKey(key)
	}
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		s.positiveFilters[key] = &valsFilter{Vals: vals}
		delete(s.negativeFilters, key)
		s.Unlock()
	}
	return s
}

// HasNotValue sets restriction for records output.
func (s *Sink) HasNotValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return s.HasNotKey(key)
	}
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		s.negativeFilters[key] = &valsFilter{Vals: vals}
		delete(s.positiveFilters, key)
		s.Unlock()
	}
	return s
}

// Int64Range sets restriction for records output.
func (s *Sink) Int64Range(key string, from, to int64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &int64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// Int64NotRange sets restriction for records output.
func (s *Sink) Int64NotRange(key string, from, to int64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.positiveFilters, key)
		s.negativeFilters[key] = &int64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// Float64Range sets restriction for records output.
func (s *Sink) Float64Range(key string, from, to float64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &float64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// Float64NotRange sets restriction for records output.
func (s *Sink) Float64NotRange(key string, from, to float64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.positiveFilters, key)
		s.negativeFilters[key] = &float64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// TimeRange sets restriction for records output.
func (s *Sink) TimeRange(key string, from, to time.Time) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &timeRangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// TimeNotRange sets restriction for records output.
func (s *Sink) TimeNotRange(key string, from, to time.Time) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.positiveFilters, key)
		s.negativeFilters[key] = &timeRangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithFilter setup custom filtering function for values for a specific key.
// Custom filter should realize Filter interface. All custom filters treated
// as positive filters. So if the filter returns true then it will be passed.
func (s *Sink) WithFilter(key string, customFilter Filter) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = customFilter
		s.Unlock()
	}
	return s
}

// Reset all filters for the keys for the output.
func (s *Sink) Reset(keys ...string) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		for _, key := range keys {
			delete(s.positiveFilters, key)
			delete(s.negativeFilters, key)
		}
		s.Unlock()
	}
	return s
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (s *Sink) Hide(keys ...string) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		for _, key := range keys {
			s.hiddenKeys[key] = true
		}
		s.Unlock()
	}
	return s
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (s *Sink) Unhide(keys ...string) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		for _, key := range keys {
			delete(s.hiddenKeys, key)
		}
		s.Unlock()
	}
	return s
}

// Stop stops writing to the output.
func (s *Sink) Stop() *Sink {
	atomic.StoreInt32(s.state, sinkStopped)
	return s
}

// Start writing to the output.
// After creation of a new sink it will paused and you need explicitly start it.
// It allows setup the filters before the sink will accepts any records.
func (s *Sink) Start() *Sink {
	atomic.StoreInt32(s.state, sinkActive)
	return s
}

// Close closes the sink. It flushes records for the sink before closing.
func (s *Sink) Close() {
	if atomic.LoadInt32(s.state) > sinkClosed {
		atomic.StoreInt32(s.state, sinkClosed)
		s.close <- struct{}{}
		collector.Lock()
		collector.count--
		collector.sinks = append(collector.sinks[0:s.id], collector.sinks[s.id+1:]...)
		collector.Unlock()
	}
}

func processSink(s *Sink) {
	var (
		record chain
		ok     bool
	)
	for {
		select {
		case record, ok = <-s.In:
			if !ok {
				return
			}
			if atomic.LoadInt32(s.state) < sinkActive {
				record.wg.Done()
				continue
			}
			s.RLock()
			var filter Filter
			for _, pair := range record.pairs {
				// Negative conditions have highest priority
				if filter, ok = s.negativeFilters[pair.Key]; ok {
					if filter.Check(pair.Key, pair.Val) {
						goto skipRecord
					}
				}
				// At last check for positive conditions
				if filter, ok = s.positiveFilters[pair.Key]; ok {
					if !filter.Check(pair.Key, pair.Val) {
						goto skipRecord
					}
				}
			}
			s.formatRecord(record.pairs)
		skipRecord:
			s.RUnlock()
			record.wg.Done()
		case <-s.close:
			s.Lock()
			s.positiveFilters = nil
			s.negativeFilters = nil
			s.hiddenKeys = nil
			s.Unlock()
			return
		}
	}
}

func (s *Sink) formatRecord(record []*Pair) {
	s.format.Begin()
	for _, pair := range record {
		if ok := s.hiddenKeys[pair.Key]; ok {
			continue
		}
		s.format.Pair(pair.Key, pair.Val, pair.Type)
	}
	s.writer.Write(s.format.Finish())
}

const flushTimeout = 3 * time.Second

func sinkRecord(rec []*Pair) {
	var wg sync.WaitGroup
	collector.RLock()
	for _, s := range collector.sinks {
		if atomic.LoadInt32(s.state) == sinkActive {
			wg.Add(1)
			s.In <- chain{&wg, rec}
		}
	}
	collector.RUnlock()
	var c = make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return
	case <-time.After(flushTimeout):
		return
	}
}
