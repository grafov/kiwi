package kiwi_test

// It was adapted from logxi package tests.

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/grafov/kiwi"
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
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.AddInt("key", 1).AddFloat64("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Debug()
		l.AddInt("key", 1).AddFloat64("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Info()
		l.AddInt("key", 1).AddFloat64("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Warn()
		l.AddInt("key", 1).AddFloat64("key2", 3.141592).AddString("key3", "string").AddBool("key4", false).Error()
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwiTypedComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.AddInt("key", 1).AddStringer("obj", testObject).Debug()
		l.AddInt("key", 1).AddStringer("obj", testObject).Info()
		l.AddInt("key", 1).AddStringer("obj", testObject).Warn()
		l.AddInt("key", 1).AddStringer("obj", testObject).Error()
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwiTypedHelpers(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.AddPairs(
			kiwi.AsInt("key", 1),
			kiwi.AsFloat64("key2", 3.141592),
			kiwi.AsString("key3", "string"),
			kiwi.AsBool("key4", false)).Debug()
		l.AddPairs(
			kiwi.AsInt("key", 1),
			kiwi.AsFloat64("key2", 3.141592),
			kiwi.AsString("key3", "string"),
			kiwi.AsBool("key4", false)).Info()
		l.AddPairs(
			kiwi.AsInt("key", 1),
			kiwi.AsFloat64("key2", 3.141592),
			kiwi.AsString("key3", "string"),
			kiwi.AsBool("key4", false)).Warn()
		l.AddPairs(
			kiwi.AsInt("key", 1),
			kiwi.AsFloat64("key2", 3.141592),
			kiwi.AsString("key3", "string"),
			kiwi.AsBool("key4", false)).Error()
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwiTypedHelpersComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.AddPairs(kiwi.AsInt("key", 1), kiwi.AsStringer("obj", testObject)).Debug()
		l.AddPairs(kiwi.AsInt("key", 1), kiwi.AsStringer("obj", testObject)).Info()
		l.AddPairs(kiwi.AsInt("key", 1), kiwi.AsStringer("obj", testObject)).Warn()
		l.AddPairs(kiwi.AsInt("key", 1), kiwi.AsStringer("obj", testObject)).Error()
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwi(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwiComplex(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	l := kiwi.New()
	l.With("_n", "bench", "_p", pid)
	l.WithTimestamp(time.RFC3339)
	kiwi.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "obj", testObject)
		l.Info("key", 1, "obj", testObject)
		l.Warn("key", 1, "obj", testObject)
		l.Error("key", 1, "obj", testObject)
	}
	b.StopTimer()
	out.Flush().Close()
}

func BenchmarkLevelsKiwiGlobal(b *testing.B) {
	buf := &bytes.Buffer{}
	b.SetBytes(2)
	b.ResetTimer()
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "debug", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "info", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "warn", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "error", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Flush().Close()
}
