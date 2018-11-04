package kiwi

// This file consists of Sink related structures and functions.
// Outputs accepts incoming log records from Loggers, check them with filters
// and write to output streams if checks passed.

/* Copyright (c) 2016-2017, Alexander I.Grafov <grafov@gmail.com>
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
	Sinks     []*Sink
	Count     int
	WaitFlush sync.WaitGroup
}

type (
	// Sink used for filtering incoming log records from all logger instances
	// and decides how to filter them. Each output wraps its own io.Writer.
	// Sink methods are safe for concurrent usage.
	Sink struct {
		id uint
		In chan []*Pair
		/*

			l1 ptr
			l2 ptr
			l3 ptr

			[]*Pair
			[]endMarks  - ссылки на концы строк в []*Pair
		*/
		close  chan struct{}
		writer io.Writer
		format Formatter
		state  *int32

		sync.RWMutex
		positiveFilters map[string]Filter
		negativeFilters map[string]Filter
		hiddenKeys      map[string]bool
	}
)

// SinkTo creates a new sink for an arbitrary number of loggers.
// There are any number of sinks may be created for saving incoming log
// records to different places.
// The sink requires explicit start with Start() before usage.
// That allows firstly setup filters before sink will really accept any records.
func SinkTo(w io.Writer, fn Formatter) *Sink {
	collector.RLock()
	for i, sink := range collector.Sinks {
		if sink.writer == w {
			collector.Sinks[i].format = fn
			collector.RUnlock()
			return collector.Sinks[i]
		}
	}
	var state = sinkStopped
	collector.RUnlock()
	sink := &Sink{
		In:              make(chan []*Pair, 16),
		close:           make(chan struct{}),
		format:          fn,
		state:           &state,
		writer:          w,
		positiveFilters: make(map[string]Filter),
		negativeFilters: make(map[string]Filter),
		hiddenKeys:      make(map[string]bool),
	}
	collector.Lock()
	sink.id = uint(collector.Count)
	collector.Sinks = append(collector.Sinks, sink)
	collector.Count++
	go processSink(sink)
	collector.Unlock()
	return sink
}

// WithKey sets restriction for records output.
// Only the records WITH any of the keys will be passed to output.
func (s *Sink) WithKey(keys ...string) *Sink {
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

// WithoutKey sets restriction for records output.
// Only the records WITHOUT any of the keys will be passed to output.
func (s *Sink) WithoutKey(keys ...string) *Sink {
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

// WithValue sets restriction for records output.
// A record passed to output if the key equal one of any of the listed values.
func (s *Sink) WithValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return s.WithKey(key)
	}
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		s.positiveFilters[key] = &valsFilter{Vals: vals}
		delete(s.negativeFilters, key)
		s.Unlock()
	}
	return s
}

// WithoutValue sets restriction for records output.
func (s *Sink) WithoutValue(key string, vals ...string) *Sink {
	if len(vals) == 0 {
		return s.WithoutKey(key)
	}
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		s.negativeFilters[key] = &valsFilter{Vals: vals}
		delete(s.positiveFilters, key)
		s.Unlock()
	}
	return s
}

// WithInt64Range sets restriction for records output.
func (s *Sink) WithInt64Range(key string, from, to int64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &int64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithoutInt64Range sets restriction for records output.
func (s *Sink) WithoutInt64Range(key string, from, to int64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.positiveFilters, key)
		s.negativeFilters[key] = &int64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithFloat64Range sets restriction for records output.
func (s *Sink) WithFloat64Range(key string, from, to float64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &float64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithoutFloat64Range sets restriction for records output.
func (s *Sink) WithoutFloat64Range(key string, from, to float64) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.positiveFilters, key)
		s.negativeFilters[key] = &float64RangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithTimeRange sets restriction for records output.
func (s *Sink) WithTimeRange(key string, from, to time.Time) *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		s.Lock()
		delete(s.negativeFilters, key)
		s.positiveFilters[key] = &timeRangeFilter{From: from, To: to}
		s.Unlock()
	}
	return s
}

// WithoutTimeRange sets restriction for records output.
func (s *Sink) WithoutTimeRange(key string, from, to time.Time) *Sink {
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
		collector.WaitFlush.Wait()
		collector.Lock()
		s.close <- struct{}{}
		collector.Count--
		collector.Sinks = append(collector.Sinks[0:s.id], collector.Sinks[s.id+1:]...)
		collector.Unlock()
	}
}

// Flush waits that all previously sent to the output records worked.
func (s *Sink) Flush() *Sink {
	if atomic.LoadInt32(s.state) > sinkClosed {
		collector.WaitFlush.Wait()
	}
	return s
}

func processSink(s *Sink) {
	var (
		record []*Pair
		ok     bool
	)
	for {
		select {
		case record, ok = <-s.In:
			if !ok {
				atomic.StoreInt32(s.state, sinkClosed)
				s.Lock()
				s.positiveFilters = nil
				s.negativeFilters = nil
				s.hiddenKeys = nil
				s.Unlock()
				return
			}
			if atomic.LoadInt32(s.state) < sinkActive {
				collector.WaitFlush.Done()
				continue
			}
			s.RLock()
			var (
				filter Filter
			)
			for _, pair := range record {
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
			s.formatRecord(record)
		skipRecord:
			s.RUnlock()
			collector.WaitFlush.Done()
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

func sinkRecord(rec []*Pair) {
	for _, s := range collector.Sinks {
		if atomic.LoadInt32(s.state) == sinkActive {
			s.In <- rec
		} else {
			collector.WaitFlush.Done()
		}
	}
	// It was locked in Log() calls.
	collector.RUnlock()
}
