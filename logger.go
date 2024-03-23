package kiwi

import (
	"fmt"
	"sync"
)

// This file consists of Logger related structures and functions.

/* Copyright (c) 2016-2024, Alexander I.Grafov <grafov@inet.name>
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
	MessageKey = "msg"
	ErrorKey   = "kiwi-error"
	InfoKey    = "kiwi-info"
)

type (
	// Logger keeps context and log record. There are many loggers
	// initialized in different places of application. Subloggers
	// are still safe for concurrent usage. See Logger.New()
	// method below for details.
	//
	// Thread safe logic was changed in v0.6: early version was
	// not thread safe for logger instances.
	Logger struct {
		c       sync.RWMutex // sync context
		context []*Pair
		parent  *Logger
		p       sync.Mutex // sync pairs
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
// the global logger. This fuction is concurrent safe.
func Fork() *Logger {
	newContext := make([]*Pair, len(context))
	global.RLock()
	copy(newContext, context)
	global.RUnlock()
	return &Logger{context: newContext}
}

// New creates a new logger instance but not copy context from the
// global logger. The method is empty now, I keep it for compatibility
// with older versions of API.
func New() *Logger {
	return new(Logger)
}

// Fork creates a new instance of the logger. It copies the context
// from the logger from the parent logger. But the values of the
// current record of the parent logger discarded.
func (l *Logger) Fork() *Logger {
	l.c.RLock()
	fork := Logger{parent: l, context: make([]*Pair, len(l.context))}
	copy(fork.context, l.context)
	l.c.RUnlock()
	return &fork
}

// New creates a new instance of the logger. It not inherited the
// context of the parent logger. The method is empty now, I keep it
// for compatibility with older versions of API.
func (l *Logger) New() *Logger {
	return new(Logger)
}

// Log is the most common method for flushing previously added key-val pairs to an output.
// After current record is flushed all pairs removed from a record except contextSrc pairs.
func (l *Logger) Log(keyVals ...interface{}) {
	l.Add(keyVals...)
	record := make([]*Pair, 0, len(l.context)+len(l.pairs))
	l.c.RLock()
	for _, p := range l.context {
		if p.Eval != nil {
			// Evaluate delayed context value here before output.
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, nil, p.Type})
		}
	}
	l.c.RUnlock()
	l.p.Lock()
	for _, p := range l.pairs {
		if p.Eval != nil {
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, nil, p.Type})
		}
	}
	l.pairs = nil
	l.p.Unlock()
	sinkRecord(record)
}

// Add a new key-value pairs to the log record.
//
// If a key already added then value will be updated. If a key already
// exists in logger' context then it will be overriden for a current
// record only. After flushing a record with Log() old context value
// will be restored.
//
// With("key1", "value1").Log("key1", "value2") -> key1:value2
//
// Another case when a key already exists for the current record. Then
// next value with for the same key will just added again.
//
// Add("key1", "value1").Log("key1", "value2") -> key1:value1 key1:value2
func (l *Logger) Add(keyVals ...interface{}) *Logger {
	var (
		key         string
		shouldBeKey = true
	)
	l.p.Lock()
	defer l.p.Unlock()
	for _, val := range keyVals {
		if shouldBeKey {
			switch v := val.(type) {
			case string:
				key = v
			case *Pair:
				l.pairs = append(l.pairs, v)
				continue
			default:
				l.pairs = append(l.pairs, toPair(ErrorKey, fmt.Sprintf("non a string type (%T) for the key (%v)", val, val)))
				continue
			}
		} else {
			l.pairs = append(l.pairs, toPair(key, val))
		}
		shouldBeKey = !shouldBeKey
	}
	if !shouldBeKey {
		l.pairs = append(l.pairs, toPair(MessageKey, key))
	}
	return l
}

// Reset logger values added after last Log() call and logger context.
// Behavior changed: early kiwi versions (before v0.6) keep context values untouched!
func (l *Logger) Reset() *Logger {
	l.c.Lock()
	l.context = nil
	l.c.Unlock()
	l.p.Lock()
	l.pairs = nil
	l.p.Unlock()
	return l
}

// Return just return back to upper logger the context and added (not
// flushed) records. So you can log them in upper logger.
func (l *Logger) Return() {
	l.c.Lock()
	for _, c := range l.context {
		l.parent.context = append(l.parent.context, c)
	}
	l.c.Unlock()
	l.p.Lock()
	for _, p := range l.pairs {
		l.parent.pairs = append(l.parent.pairs, p)
	}
	l.p.Unlock()
}
