package kiwi

import "fmt"

// This file consists of Logger related structures and functions.

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

var (
	UnpairedKey = "message"
	ErrorKey    = "kiwi-error"
)

type (
	// Logger keeps context and log record. There are many loggers initialized
	// in different places of application. Loggers are not safe for
	// concurrent usage so then you need logger for another goroutine you will need clone existing instance.
	// See Logger.New() method below for details.
	Logger struct {
		context map[string]Pair
		pairs   []*Pair
	}
	// Stringer is the same as fmt.Stringer
	Stringer interface {
		String() string
	}
	// Valuer allows log data from any custom types if they conform this interface.
	// Also types that conform fmt.Stringer can be used. But as they not have IsQuoted() check
	// they always treated as strings and displayed in quotes.
	Valuer interface {
		Stringer
		IsQuoted() bool
	}
	// Pair is key and value together. They can be used by custom
	// helpers for example for logging timestamps or something.
	Pair struct {
		Key  string
		Val  string
		Eval interface{}
		Type int
	}
)

// Fork creates a new logger instance that inherited the context from
// the global logger.
func Fork() *Logger {
	var newContext = make(map[string]Pair, len(globalContext.m)*2)
	globalContext.RLock()
	for key, pair := range globalContext.m {
		newContext[key] = pair
	}
	globalContext.RUnlock()
	return &Logger{context: newContext}
}

// New creates a new logger instance but not copy context from the
// global logger.
func New() *Logger {
	var newContext = make(map[string]Pair)
	return &Logger{context: newContext}
}

// Fork creates a new instance of the logger. It copies the context
// from the logger from the parent logger. But the values of the
// current record of the parent logger discarded.
func (l *Logger) Fork() *Logger {
	var newContext = make(map[string]Pair, len(l.context)*2)
	for key, pair := range l.context {
		newContext[key] = pair
	}
	return &Logger{context: newContext}
}

// New creates a new instance of the logger. It not inherited the
// context of the parent logger.
func (l *Logger) New() *Logger {
	var newContext = make(map[string]Pair)
	return &Logger{context: newContext}
}

// Log is the most common method for flushing previously added key-val pairs to an output.
// After current record is flushed all pairs removed from a record except contextSrc pairs.
func (l *Logger) Log(keyVals ...interface{}) {
	// 1. Log the context.
	var record = make([]*Pair, 0, len(l.context)+len(l.pairs)+len(keyVals)/2+1)
	for _, p := range l.context {
		if p.Eval != nil {
			// Evaluate delayed context value here before output.
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, p.Eval, p.Type})
		}
	}
	// 2. Log the regular key-value pairs that added before by Add() calls.
	for _, p := range l.pairs {
		if p.Eval != nil {
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, p.Eval, p.Type})
		}
	}
	// 3. Log the regular key-value pairs that come in the args.
	var (
		key          string
		shouldBeAKey = true
	)
	for _, val := range keyVals {
		if shouldBeAKey {
			switch val.(type) {
			case string:
				key = val.(string)
			case Pair:
				record = append(record, val.(*Pair))
				continue
			default:
				record = append(record, toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val)))
				key = UnpairedKey
			}
		} else {
			record = append(record, toPair(key, val))
		}
		shouldBeAKey = !shouldBeAKey
	}
	if !shouldBeAKey && key != UnpairedKey {
		record = append(record, toPair(UnpairedKey, key))
	}
	// 4. Pass the record to the collector.
	// The collector will be unlocked inside sinkRecord().
	collector.WaitFlush.Add(collector.Count)
	collector.RLock()
	go sinkRecord(record)
	l.pairs = nil
}

// Add a new key-value pairs to the log record. If a key already added then value will be
// updated. If a key already exists in a contextSrc then it will be overridden by a new
// value for a current record only. After flushing a record with Log() old context value
// will be restored.
func (l *Logger) Add(keyVals ...interface{}) *Logger {
	var (
		key          string
		shouldBeAKey = true
	)
	// key=val pairs
	for _, val := range keyVals {
		if shouldBeAKey {
			switch val.(type) {
			case string:
				key = val.(string)
			case Pair:
				l.pairs = append(l.pairs, val.(*Pair))
				continue
			default:
				l.pairs = append(l.pairs, toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val)))
				continue
			}
		} else {
			l.pairs = append(l.pairs, toPair(key, val))
		}
		shouldBeAKey = !shouldBeAKey
	}
	if !shouldBeAKey {
		l.pairs = append(l.pairs, toPair(UnpairedKey, key))
	}
	return l
}

// With defines a context for the logger. The context overrides pairs in the record.
func (l *Logger) With(keyVals ...interface{}) *Logger {
	var (
		key          string
		shouldBeAKey = true
	)
	// key=val pairs
	for _, val := range keyVals {
		if shouldBeAKey {
			switch val.(type) {
			case string:
				key = val.(string)
			case Pair:
				l.context[key] = val.(Pair)
				continue
			case []*Pair:
				for _, p := range val.([]*Pair) {
					l.context[p.Key] = *p
				}
				continue
			default:
				l.context[ErrorKey] = *toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val))
				continue
			}
		} else {
			l.context[key] = *toPair(key, val)
		}
		shouldBeAKey = !shouldBeAKey
	}
	if !shouldBeAKey {
		l.pairs = append(l.pairs, toPair(UnpairedKey, key))
	}
	return l
}

// Without drops some keys from a context for the logger.
func (l *Logger) Without(keys ...string) *Logger {
	for _, key := range keys {
		delete(l.context, key)
	}
	return l
}

// Reset logger values added after last Log() call. It keeps context untouched.
func (l *Logger) Reset() *Logger {
	l.pairs = nil
	return l
}

// ResetContext resets the context of the logger.
func (l *Logger) ResetContext() *Logger {
	l.context = make(map[string]Pair, len(l.context)*2)
	return l
}
