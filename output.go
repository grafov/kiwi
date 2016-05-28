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
		In      chan map[string]recVal
		w       io.Writer
		format  format
		paused  bool
		filters map[string]filter // TODO think about interface{} instead of string
		hidden  map[string]bool
	}
	filter struct {
		Val   string
		Flags int8
	}
	format uint8
)

const (
	mustPresentMask int8 = 0x01
	checkValueMask  int8 = 0x02
	hiddenKeyMask   int8 = 0x04 // XXX
)

const (
	Logfmt format = iota
	JSON
)

// GetWriter creates a new output for an arbitrary number of loggers.
// There are any number of outputs may be created for saving incoming log
// records to different places.
func GetWriter(w io.Writer, logFormat format) *Output {
	outputs.Lock()
	defer outputs.Unlock()
	if out, ok := outputs.w[w]; ok {
		return out
	}
	out := &Output{
		In: make(chan map[string]recVal, 1),
		w:  w, filters: make(map[string]filter),
		format: logFormat}
	outputs.w[w] = out
	go processOutput(out)
	return out
}

// With sets restriction for log records output.
// Only records that has all keys will be logged.
func (out *Output) With(keys ...string) *Output {
	var (
		filter filter
		ok     bool
	)
	out.Lock()
	for _, tag := range keys {
		if filter, ok = out.filters[tag]; ok {
			filter.Flags |= mustPresentMask
		} else {
			filter.Flags = mustPresentMask
		}
		out.filters[tag] = filter
	}
	out.Unlock()
	return out
}

// WithPairs sets restriction for log records output.
// It will compare each key and value pair from filters
// with each key and value from logged record.
func (out *Output) WithValues(keyVals ...string) *Output {
	var (
		filter filter
		key    string
		ok     bool
	)
	out.Lock()
	for i, val := range keyVals {
		if i%2 == 0 {
			key = val
			continue
		}
		if filter, ok = out.filters[key]; ok {
			filter.Flags |= mustPresentMask
		} else {
			filter.Flags = mustPresentMask
		}
		filter.Val = val
		out.filters[key] = filter
	}
	// for odd number of key-val pairs just add label without recVal
	if len(keyVals)%2 == 1 {
		if filter, ok = out.filters[key]; ok {
			filter.Flags |= mustPresentMask
		} else {
			filter.Flags = mustPresentMask
		}
		out.filters[key] = filter
	}
	out.Unlock()
	return out
}

// Without set filter for keys those should not be present in a log record.
// It will pass only records that has no one key from this set.
func (out *Output) Without(keys ...string) *Output {
	var (
		filter filter
		ok     bool
	)
	out.Lock()
	for _, tag := range keys {
		if filter, ok = out.filters[tag]; ok {
			filter.Flags = filter.Flags &^ mustPresentMask
		} else {
			filter.Flags = 0
		}
		out.filters[tag] = filter
	}
	out.Unlock()
	return out
}

// Hide keys from the output. Other keys in record will be displayed
// but not hidden keys.
func (out *Output) Hide(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		out.hidden[tag] = true
	}
	out.Unlock()
	return out
}

// Unhide previously hidden keys. They will be displayed in the output again.
func (out *Output) Unhide(keys ...string) *Output {
	out.Lock()
	for _, tag := range keys {
		delete(out.hidden, tag)
	}
	out.Unlock()
	return out
}

func (out *Output) Pause() {
	out.paused = true
}

func (out *Output) Continue() {
	out.paused = false
}

func (out *Output) Close() {
	// TODO close channel and check
}

func broadcastRecord(record map[string]recVal) {
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
		for key, filter := range out.filters {
			if _, ok := record[key]; ok {
				if filter.Flags&mustPresentMask == 0 {
					goto skipRecord
				}
			} else {
				if filter.Flags&mustPresentMask > 0 {
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
