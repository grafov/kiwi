package kiwi_test

// It was adapted from logxi package tests.

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/grafov/kiwi"
	"github.com/grafov/kiwi/level"
	"github.com/grafov/kiwi/strict"
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

func BenchmarkLevelsKiwiStrict_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseLogfmt()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Info(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Warn(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Error(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiStrictComplex_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseLogfmt()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Info(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Warn(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Error(strict.Int("key", 1), strict.Stringer("obj", testObject))
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwi_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseLogfmt()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiComplex_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseLogfmt()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "obj", testObject)
		l.Info("key", 1, "obj", testObject)
		l.Warn("key", 1, "obj", testObject)
		l.Error("key", 1, "obj", testObject)
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiGlobal_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	out := kiwi.SinkTo(buf, kiwi.UseLogfmt()).Start()
	for i := 0; i < b.N; i++ {
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "debug", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "info", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "warn", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "error", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Close()
}
