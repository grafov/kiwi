package kiwi

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Logger keep context and log record. There are many loggers initialized
// in different places of application. Logger methods are safe for
// concurrent usage.
type (
	Logger struct {
		sync.RWMutex
		contextSrc map[interface{}]interface{}
		context    map[string]recVal
		pairs      map[string]recVal
	}
	recVal struct {
		Val  string
		Type uint8
	}
)

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

// Log is a most common method for flushing previously added key-val pairs to an output.
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
	broadcastRecord(record)
}

// Add new key-recVal pairs to the log record. If a key already added then value will be
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
		l.pairs[key] = recVal{"", emptyVal}
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
		l.context[key] = recVal{"", emptyVal}
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
	l.context["timestamp"] = recVal{time.Now().Format(format), customVal}
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

// GetContext returns copy of contextSrc saved in the logger.
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
// Most of Logger operations return *Logger itself but it made for operations
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

var levelName = "level"

// UseLevelName allows to change default recVal "level" to any recVal you want.
// Set it to empty string if you want to report level without presetting any name.
func UseLevelName(name string) {
	levelName = name
}

// Err imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "error". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Err(keyVals ...interface{}) *Logger {
	return l.Add(levelName, "error").Add(keyVals...)
}

// Warn imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "warning". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Warn(keyVals ...interface{}) *Logger {
	return l.Add(levelName, "warning").Add(keyVals...)
}

// Info imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "info". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Info(keyVals ...interface{}) *Logger {
	return l.Add(levelName, "info").Add(keyVals...)
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
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		num, ok := val.(int)
		if !ok {
			switch v := reflect.ValueOf(val); v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.FormatInt(v.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return strconv.FormatUint(v.Uint(), 10)
			default:
				return fmt.Sprintf("%v", val)
			}
		}
		return strconv.Itoa(num)
	case fmt.Stringer:
		return val.(fmt.Stringer).String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toRecordValue(val interface{}) recVal {
	if val == nil {
		return recVal{"", emptyVal}
	}
	switch val.(type) {
	case bool:
		if val.(bool) {
			return recVal{"true", booleanVal}
		}
		return recVal{"false", booleanVal}
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		num, ok := val.(int)
		if !ok {
			switch v := reflect.ValueOf(val); v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return recVal{strconv.FormatInt(v.Int(), 10), integerVal}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return recVal{strconv.FormatUint(v.Uint(), 10), integerVal}
			default:
				return recVal{fmt.Sprintf("%v", val), integerVal}
			}

		}
		return recVal{strconv.Itoa(num), integerVal}
	case fmt.Stringer:
		return recVal{val.(fmt.Stringer).String(), stringVal}
	default:
		return recVal{fmt.Sprintf("%v", val), stringVal}
	}
}
