package kiwi

import (
	"bytes"
	"io"
	"sync"
)

var outputs struct {
	sync.RWMutex
	w map[io.Writer]*Output
}

type (
	// Output used for filtering incoming log records from all logger instances
	// and decides how to write them. Each Output wraps its own io.Writer.
	// Output methods are safe for concurrent usage.
	Output struct {
		sync.RWMutex
		In              chan map[string]recVal
		w               io.Writer
		format          format
		paused          bool
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
	outputs.Lock()
	defer outputs.Unlock()
	if out, ok := outputs.w[w]; ok {
		return out
	}
	out := &Output{
		In:              make(chan map[string]recVal, 1),
		w:               w,
		positiveFilters: make(map[string]filter),
		negativeFilters: make(map[string]filter),
		format:          logFormat}
	outputs.w[w] = out
	go processOutput(out)
	return out
}

// With sets restriction for log records output.
// Only records that has all keys will be logged.
func (out *Output) With(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		delete(out.negativeFilters, tag)
		out.positiveFilters[tag] = &keyFilter{Key: tag}
	}
	out.Unlock()
	return out
}

// WithValues defines arbitrary set of values for a key.
// Any of these values for a defined key must be presented
// in a log record.
func (out *Output) WithValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return out.With(key)
	}
	out.Lock()
	delete(out.negativeFilters, key)
	out.positiveFilters[key] = &valsFilter{Key: key, Vals: vals}
	out.Unlock()
	return out
}

// Without set filter for keys those should not be present in a log record.
// It will pass only records that has no one key from this set.
func (out *Output) Without(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		delete(out.negativeFilters, tag)
		out.positiveFilters[tag] = &keyFilter{Key: tag}
	}
	out.Unlock()
	return out
}

func (out *Output) WithoutValues(key string, vals ...string) *Output {
	if len(vals) == 0 {
		return out.Without(key)
	}
	out.Lock()
	delete(out.positiveFilters, key)
	out.negativeFilters[key] = &valsFilter{Key: key, Vals: vals}
	out.Unlock()
	return out
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (out *Output) Hide(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		out.hiddenKeys[tag] = true
	}
	out.Unlock()
	return out
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (out *Output) Unhide(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		delete(out.hiddenKeys, tag)
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
	// TODO close channel and check
}

// A new record passed to all outputs. Each output routine decides n
func passRecordToOutput(record map[string]recVal) {
	outputs.RLock()
	for _, out := range outputs.w {
		out.In <- record
	}
	outputs.RUnlock()
}

func processOutput(out *Output) {
	for {
		record := <-out.In
		if out.paused {
			continue
		}
		out.RLock()
		for key, val := range record {
			if filter, ok := out.negativeFilters[key]; ok {
				if filter.Check(key, val.Val) {
					goto skipRecord
				}
			}
			if filter, ok := out.positiveFilters[key]; ok {
				if !filter.Check(key, val.Val) {
					goto skipRecord
				}
			}
		}
		out.write(record)
	skipRecord:
		out.RUnlock()
	}
}

func (out *Output) write(record map[string]recVal) {
	var logLine bytes.Buffer
	out.RLock()
	for key, val := range record {
		logLine.WriteString(key)
		logLine.WriteRune('=')
		logLine.WriteString(val.Val)
		logLine.WriteRune(' ')
	}
	out.RUnlock()
	logLine.WriteRune('\n')
	logLine.WriteTo(out.w)
}
