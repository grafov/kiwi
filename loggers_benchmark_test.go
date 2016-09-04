package kiwi_test

// It was adapted from logxi package.

import (
	"bytes"
	"encoding/json"
	L "log"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/grafov/kiwi"
	"github.com/mgutz/logxi/v1"
	"gopkg.in/inconshreveable/log15.v2"
)

type M map[string]interface{}

var testObject = M{
	"foo": "bar",
	"bah": M{
		"int":      1,
		"float":    -100.23,
		"date":     "06-01-01T15:04:05-0700",
		"bool":     true,
		"nullable": nil,
	},
}

// Right way for kiwi is realize Record interface for the custom type
// that logger can't accept directly. But you can simply pass fmt.Stringer
// interface as well.
// You need Record interface if you want specify quotation rules with IsQuoted().
// Elsewere String() is enough.
func (m M) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

var pid = os.Getpid()

func toJSON(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}

// These tests write out all log levels with concurrency turned on and
// (mostly) equivalent fields.

func BenchmarkLevelsKiwiTyped(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.UseOutput(buf, kiwi.JSON)
	defer out.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddInt("key", 1).AddFloat("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Debug()
		l.AddInt("key", 1).AddFloat("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Info()
		l.AddInt("key", 1).AddFloat("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Warn()
		l.AddInt("key", 1).AddFloat("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Error()
	}
	b.StopTimer()
}

func BenchmarkLevelsKiwiTypedComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.UseOutput(buf, kiwi.JSON)
	defer out.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddInt("key", 1).AddStringer("obj", testObject).Debug()
		l.AddInt("key", 1).AddStringer("obj", testObject).Info()
		l.AddInt("key", 1).AddStringer("obj", testObject).Warn()
		l.AddInt("key", 1).AddStringer("obj", testObject).Error()
	}
	b.StopTimer()
}

func BenchmarkLevelsKiwi(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.UseOutput(buf, kiwi.JSON)
	defer out.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}

func BenchmarkLevelsKiwiComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.UseOutput(buf, kiwi.JSON)
	defer out.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "obj", testObject)
		l.Info("key", 1, "obj", testObject)
		l.Warn("key", 1, "obj", testObject)
		l.Error("key", 1, "obj", testObject)
	}
	b.StopTimer()
}

func BenchmarkLevelsStdLog(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := L.New(buf, "bench ", L.LstdFlags)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		debug := map[string]interface{}{"l": "debug", "key1": 1, "key2": 3.141592, "key3": "string", "key4": false}
		l.Printf(toJSON(debug))

		info := map[string]interface{}{"l": "info", "key1": 1, "key2": 3.141592, "key3": "string", "key4": false}
		l.Printf(toJSON(info))

		warn := map[string]interface{}{"l": "warn", "key1": 1, "key2": 3.141592, "key3": "string", "key4": false}
		l.Printf(toJSON(warn))

		err := map[string]interface{}{"l": "error", "key1": 1, "key2": 3.141592, "key3": "string", "key4": false}
		l.Printf(toJSON(err))
	}
	b.StopTimer()
}

func BenchmarkLevelsStdLogComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := L.New(buf, "bench ", L.LstdFlags)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		debug := map[string]interface{}{"l": "debug", "key1": 1, "obj": testObject}
		l.Printf(toJSON(debug))

		info := map[string]interface{}{"l": "info", "key1": 1, "obj": testObject}
		l.Printf(toJSON(info))

		warn := map[string]interface{}{"l": "warn", "key1": 1, "obj": testObject}
		l.Printf(toJSON(warn))

		err := map[string]interface{}{"l": "error", "key1": 1, "obj": testObject}
		l.Printf(toJSON(err))
	}
	b.StopTimer()
}

func BenchmarkLevelsLogxi(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	stdout := log.NewConcurrentWriter(buf)
	l := log.NewLogger3(stdout, "bench", log.NewJSONFormatter("bench"))
	l.SetLevel(log.LevelDebug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("info", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("warn", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("error", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}

func BenchmarkLevelsLogxiComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	stdout := log.NewConcurrentWriter(buf)
	l := log.NewLogger3(stdout, "bench", log.NewJSONFormatter("bench"))
	l.SetLevel(log.LevelDebug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "obj", testObject)
		l.Info("info", "key", 1, "obj", testObject)
		l.Warn("warn", "key", 1, "obj", testObject)
		l.Error("error", "key", 1, "obj", testObject)
	}
	b.StopTimer()

}

func BenchmarkLevelsLogrus(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := logrus.New()
	l.Out = buf
	l.Formatter = &logrus.JSONFormatter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "key2": 3.141592, "key3": "string", "key4": false}).Debug("debug")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "key2": 3.141592, "key3": "string", "key4": false}).Info("info")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "key2": 3.141592, "key3": "string", "key4": false}).Warn("warn")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "key2": 3.141592, "key3": "string", "key4": false}).Error("error")
	}
	b.StopTimer()
}

func BenchmarkLevelsLogrusComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := logrus.New()
	l.Out = buf
	l.Formatter = &logrus.JSONFormatter{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "obj": testObject}).Debug("debug")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "obj": testObject}).Info("info")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "obj": testObject}).Warn("warn")
		l.WithFields(logrus.Fields{"_n": "bench", "_p": pid, "key": 1, "obj": testObject}).Error("error")
	}
	b.StopTimer()
}

func BenchmarkLevelsLog15(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := log15.New(log15.Ctx{"_n": "bench", "_p": pid})
	l.SetHandler(log15.SyncHandler(log15.StreamHandler(buf, log15.JsonFormat())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("info", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("warn", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("error", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()

}

func BenchmarkLevelsLog15Complex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := log15.New(log15.Ctx{"_n": "bench", "_p": pid})
	l.SetHandler(log15.SyncHandler(log15.StreamHandler(buf, log15.JsonFormat())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "obj", testObject)
		l.Info("info", "key", 1, "obj", testObject)
		l.Warn("warn", "key", 1, "obj", testObject)
		l.Error("error", "key", 1, "obj", testObject)
	}
	b.StopTimer()
}

/*
$ go test -bench=. -benchmem
BenchmarkLevelsKiwi-4            	   50000	     31437 ns/op	   0.06 MB/s	    6993 B/op	     115 allocs/op
BenchmarkLevelsKiwiComplex-4     	   30000	     57841 ns/op	   0.03 MB/s	   10303 B/op	     173 allocs/op
BenchmarkLevelsStdLog-4          	  100000	     21769 ns/op	   0.09 MB/s	    7159 B/op	     124 allocs/op
BenchmarkLevelsStdLogComplex-4   	   50000	     35168 ns/op	   0.06 MB/s	   11446 B/op	     200 allocs/op
BenchmarkLevelsLogxi-4           	  100000	     13628 ns/op	   0.15 MB/s	    4127 B/op	      74 allocs/op
BenchmarkLevelsLogxiComplex-4    	   50000	     31377 ns/op	   0.06 MB/s	    8713 B/op	     162 allocs/op
BenchmarkLevelsLogrus-4          	   50000	     38582 ns/op	   0.05 MB/s	   12320 B/op	     177 allocs/op
BenchmarkLevelsLogrusComplex-4   	   30000	     45749 ns/op	   0.04 MB/s	   13990 B/op	     231 allocs/op
BenchmarkLevelsLog15-4           	   30000	     53837 ns/op	   0.04 MB/s	   14998 B/op	     224 allocs/op
BenchmarkLevelsLog15Complex-4    	   20000	     61730 ns/op	   0.03 MB/s	   15127 B/op	     246 allocs/op
*/
