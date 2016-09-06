package kiwi

/*
Copyright (c) 2016, Alexander I.Grafov aka Axel
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

ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"time"
)

var outputs []*Output

type (
	// Output used for filtering incoming log records from all logger instances
	// and decides how to filter them. Each output wraps its own io.Writer.
	// Output methods are safe for concurrent usage.
	Output struct {
		sync.RWMutex
		In              chan *[]pair
		W               io.Writer
		format          format
		paused          bool
		closed          bool
		positiveFilters map[string]filter
		negativeFilters map[string]filter
		hiddenKeys      map[string]bool
	}
	// Formatter hasn't used yet. Just added for future realization.
	Formatter interface {
		Format(*Output)
	}
	format uint8
)

// current output formats
const (
	Logfmt format = iota
	JSON
)

// UseOutput creates a new output for an arbitrary number of loggers.
// There are any number of outputs may be created for saving incoming log
// records to different places.
func UseOutput(w io.Writer, logFormat format) *Output {
	for _, output := range outputs {
		if output.W == w {
			return output
		}
	}
	output := &Output{
		In:              make(chan *[]pair, 16),
		W:               w,
		positiveFilters: make(map[string]filter),
		negativeFilters: make(map[string]filter),
		format:          logFormat}
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

// XXX fix it
func (o *Output) Flush() {
	o.In <- nil                      // XXX
	time.Sleep(5 * time.Millisecond) // XXX
}

// A new record passed to all outputs. Each output routine decides n
func passRecordToOutput(record []pair) {
	for _, o := range outputs {
		if !o.closed && !o.paused {
			o.In <- &record
		}
	}
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
		if record == nil {
			// Flush!
			time.Sleep(5 * time.Millisecond) // XXX make real flush
			continue
		}
		o.RLock()
		for _, pair := range *record {
			if filter, ok := o.negativeFilters[pair.Key]; ok {
				if filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
			if filter, ok := o.positiveFilters[pair.Key]; ok {
				if !filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
		}
		o.RUnlock()
		o.filter(record)
		continue
	skipRecord:
		o.RUnlock()
	}
}

// TODO separate filter from formatter
func (o *Output) filter(record *[]pair) {
	var logLine bytes.Buffer
	switch o.format {
	case JSON:
		logLine.WriteRune('{')
		o.RLock()
		for _, pair := range *record {
			if ok := o.hiddenKeys[pair.Key]; ok {
				continue
			}
			logLine.WriteRune('"')
			logLine.WriteString(pair.Key)
			logLine.WriteString("\":")
			if pair.Val.Quoted {
				logLine.WriteString(strconv.Quote(pair.Val.Strv))
			} else {
				logLine.WriteString(pair.Val.Strv)
			}
			logLine.WriteString(", ")
		}
		logLine.WriteRune('}')
	case Logfmt:
		o.RLock()
		for _, pair := range *record {
			if ok := o.hiddenKeys[pair.Key]; ok {
				continue
			}
			logLine.WriteString(pair.Key)
			logLine.WriteRune('=')
			if pair.Val.Quoted {
				logLine.WriteString(strconv.Quote(pair.Val.Strv))
			} else {
				logLine.WriteString(pair.Val.Strv)
			}
			logLine.WriteRune(' ')
		}
	}
	o.RUnlock()
	logLine.WriteRune('\n')
	logLine.WriteTo(o.W)
}
