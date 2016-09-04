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
	// and decides how to write them. Each Output wraps its own io.Writer.
	// Output methods are safe for concurrent usage.
	Output struct {
		sync.RWMutex
		//		In              chan map[string]value // TODO список вместо мапы
		In              chan *[]pair
		w               io.Writer
		format          format
		paused          bool
		closed          bool
		positiveFilters map[string]filter
		negativeFilters map[string]filter
		hiddenKeys      map[string]bool
	}
	format uint8
)

// const (
// 	mustPresentMask int8 = 0x01
// 	checkValueMask  int8 = 0x02
// )

const (
	Logfmt format = iota
	JSON
)

// UseOutput creates a new output for an arbitrary number of loggers.
// There are any number of outputs may be created for saving incoming log
// records to different places.
func UseOutput(w io.Writer, logFormat format) *Output {
	for _, out := range outputs {
		if out.w == w {
			return out
		}
	}
	out := &Output{
		In:              make(chan *[]pair, 16),
		w:               w,
		positiveFilters: make(map[string]filter),
		negativeFilters: make(map[string]filter),
		format:          logFormat}
	outputs = append(outputs, out)
	go processOutput(out)
	return out
}

// With sets restriction for records output.
// Only the records WITH any of the keys will be passed to output.
func (out *Output) With(keys ...string) *Output {
	out.Lock()
	if !out.closed {
		for _, tag := range keys {
			delete(out.negativeFilters, tag)
			out.positiveFilters[tag] = &keyFilter{Key: tag}
		}
	}
	out.Unlock()
	return out
}

// Without sets restriction for records output.
// Only the records WITHOUT any of the keys will be passed to output.
func (out *Output) Without(keys ...string) *Output {
	out.Lock()
	if !out.closed {
		for _, tag := range keys {
			out.negativeFilters[tag] = &keyFilter{Key: tag}
			delete(out.positiveFilters, tag)
		}
	}
	out.Unlock()
	return out
}

// WithValues sets restriction for records output.
// A record passed to output if the key equal one of any of the listed values.
func (out *Output) WithValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return out.With(key)
	}
	out.Lock()
	if !out.closed {
		delete(out.negativeFilters, key)
		out.positiveFilters[key] = &valsFilter{Key: key, Vals: vals}
	}
	out.Unlock()
	return out
}

// WithoutValues sets restriction for records output.
func (out *Output) WithoutValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return out.Without(key)
	}
	out.Lock()
	if !out.closed {
		delete(out.positiveFilters, key)
		out.negativeFilters[key] = &valsFilter{Key: key, Vals: vals}
	}
	out.Unlock()
	return out
}

// WithRangeInt64 sets restriction for records output.
func (out *Output) WithRangeInt64(key string, from, to int64) *Output {
	out.Lock()
	if !out.closed {
		delete(out.negativeFilters, key)
		out.positiveFilters[key] = &rangeInt64Filter{Key: key, From: from, To: to}
	}
	out.Unlock()
	return out
}

// WithRangeFloat64 sets restriction for records output.
func (out *Output) WithoutRangeInt64(key string, from, to int64) *Output {
	return out
}

// WithRangeFloat64 sets restriction for records output.
func (out *Output) WithRangeFloat64(key string, from, to float64) *Output {
	out.Lock()
	delete(out.negativeFilters, key)
	out.positiveFilters[key] = &rangeFloat64Filter{Key: key, From: from, To: to}
	out.Unlock()
	return out
}

// WithoutRangeFloat64  sets restriction for records output.
func (out *Output) WithoutRangeFloat64(key string, from, to float64) *Output {
	// XXX
	return out
}

// Reset all filters for the keys for the output.
func (out *Output) Reset(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		delete(out.positiveFilters, tag)
		delete(out.negativeFilters, tag)
	}
	out.Unlock()
	return out
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (out *Output) Hide(keys ...string) *Output {
	out.Lock()
	if !out.closed {
		for _, tag := range keys {
			out.hiddenKeys[tag] = true
		}
	}
	out.Unlock()
	return out
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (out *Output) Unhide(keys ...string) *Output {
	out.Lock()
	if !out.closed {
		for _, tag := range keys {
			delete(out.hiddenKeys, tag)
		}
	}
	out.Unlock()
	return out
}

// Pause stops writing to the output.
func (out *Output) Pause() {
	out.paused = true
}

// Contiunue writing to the output.
func (out *Output) Continue() {
	out.paused = false
}

func (out *Output) Close() {
	out.Lock()
	out.closed = true
	out.Unlock()
	close(out.In)
}

func (out *Output) Flush() {
	out.In <- nil
	time.Sleep(5 * time.Millisecond)
}

// A new record passed to all outputs. Each output routine decides n
func passRecordToOutput(record []pair) {
	for _, out := range outputs {
		if !out.closed && !out.paused {
			out.In <- &record
		}
	}
}

func processOutput(out *Output) {
	for {
		record, ok := <-out.In
		if !ok {
			out.positiveFilters = nil
			out.negativeFilters = nil
			out.hiddenKeys = nil
			return
		}
		if out.closed || out.paused {
			continue
		}
		if record == nil {
			// Flush!
			time.Sleep(5 * time.Millisecond) // XXX make real flush
			continue
		}
		out.RLock()
		for _, pair := range *record {
			if filter, ok := out.negativeFilters[pair.Key]; ok {
				if filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
			if filter, ok := out.positiveFilters[pair.Key]; ok {
				if !filter.Check(pair.Key, pair.Val.Strv) {
					goto skipRecord
				}
			}
		}
		out.RUnlock()
		out.write(record)
		continue
	skipRecord:
		out.RUnlock()
	}
}

// it yet ignores output format
func (out *Output) write(record *[]pair) {
	var logLine bytes.Buffer
	switch out.format {
	case JSON:
		logLine.WriteRune('{')
		out.RLock()
		for _, pair := range *record {
			if ok := out.hiddenKeys[pair.Key]; ok {
				continue
			}
			logLine.WriteRune('"')
			logLine.WriteString(pair.Key)
			logLine.WriteString("\":")
			var curVal string
			if pair.Val.Func != nil {
				// Evaluate lazy value here
				// TODO exaluation should be made in logger not output!
				// XXX refactor it!
				tmp := toRecordValue(toFunc(pair.Val.Func))
				curVal = tmp.Strv
			} else {
				curVal = pair.Val.Strv
			}
			if pair.Val.Quoted {
				logLine.WriteString(strconv.Quote(curVal))
			} else {
				logLine.WriteString(curVal)
			}
			logLine.WriteString(", ")
		}
		logLine.WriteRune('}')
	case Logfmt:
		out.RLock()
		for _, pair := range *record {
			if ok := out.hiddenKeys[pair.Key]; ok {
				continue
			}
			logLine.WriteString(pair.Key)
			logLine.WriteRune('=')
			var curVal string
			if pair.Val.Func != nil {
				// Evaluate lazy value here
				// TODO exaluation should be made in logger not output!
				// XXX refactor it!
				tmp := toRecordValue(toFunc(pair.Val.Func))
				curVal = tmp.Strv
			} else {
				curVal = pair.Val.Strv
			}
			if pair.Val.Quoted {
				logLine.WriteString(strconv.Quote(curVal))
			} else {
				logLine.WriteString(curVal)
			}
			logLine.WriteRune(' ')
		}
	}
	out.RUnlock()
	logLine.WriteRune('\n')
	logLine.WriteTo(out.w)
}
