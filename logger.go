package kiwi

import (
	"sync"
	"time"
)

type (
	// Logger keeps context and log record. There are many loggers initialized
	// in different places of application. Loggers are safe for
	// concurrent usage.
	Logger struct {
		sync.RWMutex
		contextSrc map[interface{}]interface{}
		context    map[string]value
		pairs      map[string]value
	}
	// Record allows log data from any custom types in they conform this interface.
	// Also types that conform fmt.Stringer can be used. But as they not have IsQuoted() check
	// they always treated as strings and displayed in quotes.
	Record interface {
		String() string
		IsQuoted() bool
	}
	value struct {
		Val    string
		Func   interface{}
		Type   uint8
		Quoted bool
	}
)

// NewLogger creates logger instance.
func NewLogger() *Logger {
	return &Logger{
		contextSrc: make(map[interface{}]interface{}),
		context:    make(map[string]value),
		pairs:      make(map[string]value)}
}

// Copy creates copy of logger instance.
func (l *Logger) Copy() *Logger {
	var (
		newContextSrc = make(map[interface{}]interface{})
		newContext    = make(map[string]value)
	)
	l.RLock()
	for k, v := range l.contextSrc {
		newContextSrc[k] = v
	}
	for k, v := range l.context {
		newContext[k] = v
	}
	l.RUnlock()
	return &Logger{
		contextSrc: newContextSrc,
		context:    newContext,
		pairs:      make(map[string]value)}
}

// Log is the most common method for flushing previously added key-val pairs to an output.
// After current record is flushed all pairs removed from a record except contextSrc pairs.
func (l *Logger) Log(keyVals ...interface{}) {
	if len(keyVals) > 0 {
		l.Add(keyVals...)
	}
	l.Lock()
	record := l.pairs
	l.pairs = make(map[string]value)
	l.Unlock()
	var key string
	for i, val := range keyVals {
		if i%2 == 0 {
			key = toRecordKey(val)
			continue
		}
		record[key] = toRecordValue(val)
	}
	// for odd number of key-val pairs just add label without value
	if len(keyVals)%2 == 1 {
		record[key] = value{"", nil, voidVal, false}
	}
	l.RLock()
	for key, val := range l.context {
		// pairs override context
		if _, ok := record[key]; !ok {
			record[key] = val
		}
	}
	l.RUnlock()
	passRecordToOutput(record)
}

// Add a new key-value pairs to the log record. If a key already added then value will be
// updated. If a key already exists in a contextSrc then it will be overriden by a new
// value for a current record only. After flushing a record with Log() old context value
// will be restored.
func (l *Logger) Add(keyVals ...interface{}) *Logger {
	var key string
	l.Lock()
	for i, val := range keyVals {
		if i%2 == 0 {
			key = toRecordKey(val)
			continue
		}
		l.pairs[key] = toRecordValue(val)
	}
	// for odd number of key-val pairs just add label without value
	if len(keyVals)%2 == 1 {
		l.pairs[key] = value{"", nil, voidVal, false}
	}
	l.Unlock()
	return l
}

// With defines a context for the logger.
func (l *Logger) With(keyVals ...interface{}) *Logger {
	var (
		keySrc interface{}
		key    string
	)
	l.Lock()
	for i, val := range keyVals {
		if i%2 == 0 {
			keySrc = val
			key = toRecordKey(val)
			continue
		}
		l.contextSrc[keySrc] = val
		l.context[key] = toRecordValue(val)
	}
	// for odd number of key-val pairs just add label without value
	if len(keyVals)%2 == 1 {
		l.contextSrc[keySrc] = nil
		l.context[key] = value{"", nil, voidVal, false}
	}
	l.Unlock()
	return l
}

// Without drops some keys from a context for the logger.
func (l *Logger) Without(keys ...interface{}) *Logger {
	l.Lock()
	for _, key := range keys {
		if _, ok := l.contextSrc[key]; ok {
			delete(l.contextSrc, key)
			delete(l.context, toRecordKey(key))
		}
	}
	l.Unlock()
	return l
}

// WithTimestamp adds "timestamp" field to the context.
func (l *Logger) WithTimestamp(format string) *Logger {
	l.Lock()
	l.contextSrc["timestamp"] = func() string { return time.Now().Format(format) }
	l.context["timestamp"] = value{"", func() string { return time.Now().Format(format) }, stringVal, true}
	l.Unlock()
	return l
}

// Reset logger values added after last Log() call. It keeps contextSrc untouched.
func (l *Logger) Reset() *Logger {
	l.Lock()
	l.pairs = make(map[string]value)
	l.Unlock()
	return l
}

// ResetContext resets the context of the logger.
func (l *Logger) ResetContext() *Logger {
	l.Lock()
	l.contextSrc = make(map[interface{}]interface{})
	l.context = make(map[string]value)
	l.Unlock()
	return l
}

// GetContext returns copy of the context saved in the logger.
func (l *Logger) GetContext() map[interface{}]interface{} {
	var contextSrc = make(map[interface{}]interface{})
	l.RLock()
	for k, v := range l.contextSrc {
		contextSrc[k] = v
	}
	l.RUnlock()
	return contextSrc
}

// GetContextValue returns single context value for the key.
func (l *Logger) GetContextValue(key string) interface{} {
	l.RLock()
	value := l.contextSrc[key]
	l.RUnlock()
	return value
}

// GetRecord returns copy of current set of keys and values prepared for logging
// as strings. With context key-vals included.
// The most of Logger operations return *Logger itself but it made for operations
// chaining only. If you need get log pairs use GelRecord() for it.
func (l *Logger) GetRecord() map[string]string {
	var merged = make(map[string]string)
	l.RLock()
	for k, v := range l.context {
		merged[k] = v.Val
	}
	for k, v := range l.pairs {
		merged[k] = v.Val
	}
	l.RUnlock()
	return merged
}

// Flush confirms that all outputs got the last logged record.
func (l *Logger) Flush() {
	// XXX
	time.Sleep(100 * time.Millisecond)
}
