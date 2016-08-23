package kiwi

import (
	"fmt"
	"strconv"
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
		context    map[string]recVal
		pairs      map[string]recVal
	}
	// Record allows log data from any custom types in they conform this interface.
	// Also types that conform fmt.Stringer can be used. But as they not have IsQuoted() check
	// they always treated as strings and displayed in quotes.
	Record interface {
		String() string
		IsQuoted() bool
	}
	recVal struct {
		Val    string
		Type   uint8
		Quoted bool
	}
)

// obsoleted by recVal iface
const (
	emptyVal uint8 = iota
	stringVal
	integerVal
	floatVal
	booleanVal
	customVal // for use with `func() string`
)

// NewLogger creates logger instance.
func NewLogger() *Logger {
	return &Logger{
		contextSrc: make(map[interface{}]interface{}),
		context:    make(map[string]recVal),
		pairs:      make(map[string]recVal)}
}

// Log is the most common method for flushing previously added key-val pairs to an output.
// After current record is flushed all pairs removed from a record except contextSrc pairs.
func (l *Logger) Log(keyVals ...interface{}) {
	if len(keyVals) > 0 {
		l.Add(keyVals...)
	}
	l.Lock()

	record := l.pairs
	l.pairs = make(map[string]recVal)
	for key, val := range l.context {
		// pairs override context
		if _, ok := record[key]; !ok {
			record[key] = val
		}
	}
	l.Unlock()
	passRecordToOutput(record)
}

// Add a new key-recVal pairs to the log record. If a key already added then value will be
// updated. If a key already exists in a contextSrc then it will be overriden by a new
// recVal for a current record only. After flushing a record with Log() old context value
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
	// for odd number of key-val pairs just add label without recVal
	if len(keyVals)%2 == 1 {
		l.pairs[key] = recVal{"", emptyVal, false}
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
	// for odd number of key-val pairs just add label without recVal
	if len(keyVals)%2 == 1 {
		l.contextSrc[keySrc] = nil
		l.context[key] = recVal{"", emptyVal, false}
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

// WithTimestamp adds "timestamp" field to a context.
func (l *Logger) WithTimestamp(format string) *Logger {
	l.Lock()
	// TODO think about offer fmt.Stringer here instead of custom func?
	l.contextSrc["timestamp"] = func() string { return time.Now().Format(format) }
	l.context["timestamp"] = recVal{time.Now().Format(format), customVal, true}
	l.Unlock()
	return l
}

// Reset logger values added after last Log() call. It keeps contextSrc untouched.
func (l *Logger) Reset() *Logger {
	l.Lock()
	l.pairs = make(map[string]recVal)
	l.Unlock()
	return l
}

// ResetContext resets the context of the logger.
func (l *Logger) ResetContext() *Logger {
	l.Lock()
	l.contextSrc = make(map[interface{}]interface{})
	l.context = make(map[string]recVal)
	l.Unlock()
	return l
}

// GetContext returns copy of context saved in the logger.
func (l *Logger) GetContext() map[interface{}]interface{} {
	var contextSrc = make(map[interface{}]interface{})
	l.RLock()
	for k, v := range l.contextSrc {
		contextSrc[k] = v
	}
	l.RUnlock()
	return contextSrc
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

func toRecordKey(val interface{}) string {
	if val == nil {
		return ""
	}
	switch val.(type) {
	case bool:
		if val.(bool) {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(val.(int))
	case int8:
		return strconv.FormatInt(int64(val.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(val.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(val.(int64)), 10)
	case uint8:
		return strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(val.(uint64)), 10)
	case fmt.Stringer:
		return val.(fmt.Stringer).String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toRecordValue(val interface{}) recVal {
	if val == nil {
		return recVal{"", emptyVal, false}
	}
	switch val.(type) {
	case bool:
		if val.(bool) {
			return recVal{"true", booleanVal, false}
		}
		return recVal{"false", booleanVal, false}
	case int:
		return recVal{strconv.Itoa(val.(int)), integerVal, false}
	case int8:
		return recVal{strconv.FormatInt(int64(val.(int8)), 10), integerVal, false}
	case int16:
		return recVal{strconv.FormatInt(int64(val.(int16)), 10), integerVal, false}
	case int32:
		return recVal{strconv.FormatInt(int64(val.(int32)), 10), integerVal, false}
	case int64:
		return recVal{strconv.FormatInt(int64(val.(int64)), 10), integerVal, false}
	case uint8:
		return recVal{strconv.FormatUint(uint64(val.(uint8)), 10), integerVal, false}
	case uint16:
		return recVal{strconv.FormatUint(uint64(val.(uint16)), 10), integerVal, false}
	case uint32:
		return recVal{strconv.FormatUint(uint64(val.(uint32)), 10), integerVal, false}
	case uint64:
		return recVal{strconv.FormatUint(uint64(val.(uint64)), 10), integerVal, false}
	case Record:
		return recVal{val.(Record).String(), stringVal, true}
	case fmt.Stringer:
		return recVal{val.(fmt.Stringer).String(), stringVal, true}
	default:
		return recVal{fmt.Sprintf("%v", val), stringVal, true}
	}
}
