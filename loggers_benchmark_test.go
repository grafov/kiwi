package kiwi_test

// It was adapted from logxi package.

import (
	"encoding/json"
	L "log"
	"os"
	"testing"
	"time"

	"bytes"
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
func (m *M) String() string {
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

func BenchmarkLog(b *testing.B) {
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

func BenchmarkLogComplex(b *testing.B) {
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

func BenchmarkKiwi(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.NewLogger()
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

func BenchmarkKiwiComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	l := kiwi.NewLogger()
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

func BenchmarkLogxi(b *testing.B) {
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

func BenchmarkLogxiComplex(b *testing.B) {
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

func BenchmarkLogrus(b *testing.B) {
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

func BenchmarkLogrusComplex(b *testing.B) {
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

func BenchmarkLog15(b *testing.B) {
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

func BenchmarkLog15Complex(b *testing.B) {
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
