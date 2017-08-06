package kiwi

// This file consists of Logger related structures and functions.

/* Copyright (c) 2016-2017, Alexander I.Grafov
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

// New creates a new logger instance.
func New() *Logger {
	return &Logger{context: make(map[string]Pair)}
}

// New creates copy of logger instance. It copies the context of the old logger
// but skips values of the current record of the old logger.
func (l *Logger) New() *Logger {
	var newContext = make(map[string]Pair, len(l.context))
	for key, pair := range l.context {
		newContext[key] = pair
	}
	return &Logger{context: newContext}
}

// Log is the most common method for flushing previously added key-val pairs to an output.
// After current record is flushed all pairs removed from a record except contextSrc pairs.
func (l *Logger) Log(keyVals ...interface{}) {
	var record = make([]*Pair, 0, len(l.context)+len(l.pairs)+len(keyVals)/2+1)
	for _, p := range l.context {
		if p.Eval != nil {
			// Evaluate delayed context value here before output.
			record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
		} else {
			record = append(record, &Pair{p.Key, p.Val, p.Eval, p.Type})
		}
	}
	for _, p := range l.pairs {
		if p.Type != deleted {
			if p.Eval != nil {
				record = append(record, &Pair{p.Key, p.Eval.(func() string)(), p.Eval, p.Type})
			} else {
				record = append(record, &Pair{p.Key, p.Val, p.Eval, p.Type})
			}
		}
	}
	var (
		key     string
		nextKey = true
	)
	for _, val := range keyVals {
		if nextKey {
			switch val.(type) {
			case Pair:
				l.pairs = append(l.pairs, val.(*Pair))
				continue
			}
			key = toKey(val)
			nextKey = false
		} else {
			var p *Pair
			if p = toPair(key, val); p.Eval != nil {
				p.Val = p.Eval.(func() string)()
			}
			record = append(record, p)
			nextKey = true
		}

	}
	// TODO for odd number of arguments pass the last argument as a
	// value with some predefined key ("message" for
	// example). Usecase: Log("Just a message without a key")
	if !nextKey {
		record = append(record, &Pair{key, "", nil, VoidVal})
	}
	collector.WaitFlush.Add(collector.Count)
	// It will be unlocked inside sinkRecord().
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
		key     string
		nextKey = true
	)
	// key=val pairs
	for _, val := range keyVals {
		if nextKey {
			switch val.(type) {
			case Pair:
				l.pairs = append(l.pairs, val.(*Pair))
				continue
			}
			key = toKey(val)
			nextKey = false
		} else {
			l.pairs = append(l.pairs, toPair(key, val))
			nextKey = true
		}
	}
	//  add a key without value for odd number for key-val pairs
	if !nextKey {
		l.pairs = append(l.pairs, &Pair{key, "", nil, VoidVal})
	}
	return l
}

// With defines a context for the logger. The context overrides pairs in the record.
func (l *Logger) With(keyVals ...interface{}) *Logger {
	var (
		key     string
		nextKey = true
	)
	// key=val pairs
	for _, val := range keyVals {
		if nextKey {
			switch val.(type) {
			case Pair:
				l.context[key] = val.(Pair)
				continue
			}
			key = toKey(val)
			nextKey = false
			continue
		}
		l.context[key] = *toPair(key, val)
		nextKey = true
	}
	// add a key without value for odd number for key-val pairs
	if !nextKey {
		l.context[key] = Pair{key, "", nil, VoidVal}
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
	l.context = make(map[string]Pair, len(l.context))
	return l
}
