package kiwi

import (
	"bytes"
	//	"reflect"
	"strings"
	"testing"
)

/* All tests consists of three parts:

- arrange structures and initialize objects for use in tests
- act on testing object
- check and assert on results

These parts separated by empty line in each test function.
*/

var (
	sampleContext = []interface{}{"context1", "value", "context2", 1, "common", []string{"the", "same"}}
	sampleRecord  = []interface{}{"key1", "value", "key2", 2, 3, 4, "common", []string{"the", "same"}}
)

// Get records from logger. Helper for testing.
func (l *Logger) getRecords() []pair {
	return l.pairs
}

// Get context from logger. Helper for testing.
func (l *Logger) getContext() []pair {
	return l.context
}

func TestNewLogger(t *testing.T) {
	l := New()

	if l == nil {
		t.Fatal("initalized logger is nil")
	}
}

// XXX
func TestLogger_With(t *testing.T) {

}

// func TestLogger_Add(t *testing.T) {
// 	l := New()

// 	l.Add(sampleRecord...)

// 	records := l.getRecords()
// 	var key string
// 	for i, sampleVal := range sampleRecord {
// 		if i%2 == 0 {
// 			key = toRecordKey(sampleVal)
// 			continue
// 		}
// 		if savedVal, ok := records[key]; ok {
// 			if reflect.DeepEqual(savedVal, sampleVal) {
// 				t.Fatalf("values not equal %v %v", savedVal, sampleVal)
// 			}
// 		} else {
// 			t.Fatalf("key %v not found", key)
// 		}
// 	}
// }

// // //
// // func TestLogger_Get_RecordsOnly(t *testing.T) {
// // 	l := New()
// // 	l.Add(sampleRecord...)

// // 	records := l.GetRecord()
// // 	for key, sampleV := range sampleRecord {
// // 		if savedVal, ok := records[key]; ok {
// // 			if savedVal != sampleV {
// // 				t.Fatalf("values not equal %v %v", savedVal, sampleV)
// // 			}
// // 		} else {
// // 			t.Fatalf("key %v not found", key)
// // 		}
// // 	}
// // }

// // // XXX
// // func TestLogger_GetLog_ContextOnly(t *testing.T) {
// // 	l := New()

// // 	l.With(sampleContext...)

// // 	context := l.GetRecord()
// // 	for key, sampleV := range sampleContext {
// // 		if savedVal, ok := context[key]; ok {
// // 			if savedVal != sampleV {
// // 				t.Fatalf("values not equal %v %v", savedVal, sampleV)
// // 			}
// // 		} else {
// // 			t.Fatalf("key %v not found", key)
// // 		}
// // 	}
// // }

// // // XXX
// // func TestLogger_GetLog_ContextOverridenByRecords(t *testing.T) {
// // 	l := New()

// // 	l.With(sampleContext...).Add(sampleRecord...)

// // 	records := l.getRecords()
// // 	// XXX context := l.getContext()
// // 	for key, sampleV := range sampleRecord {
// // 		if savedVal, ok := records[key]; ok {
// // 			if savedVal != sampleV {
// // 				t.Fatalf("values not equal %v %v", savedVal, sampleV)
// // 			}
// // 		} else {
// // 			t.Fatalf("key %v not found", key)
// // 		}
// // 	}
// // }

// // func TestLogger_Reset(t *testing.T) {
// // 	l := New()
// // 	l.Add(sampleRecord...)

// // 	l.Reset()

// // 	if len(l.GetRecord()) > 0 {
// // 		t.Fatal("reset doesn't works")
// // 	}
// // }

func TestLogger_Add_Chained(t *testing.T) {
	log := New().With(sampleContext...).Add(sampleRecord...)

	log.Log()
	log.Add("key", "value2").Log()
}

func TestLogger_IntValues(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := UseOutput(output, Logfmt)
	defer out.Close()

	log.Log("k", 123)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=123" { // XXX
		t.Fail()
	}
}
