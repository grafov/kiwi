package kiwi_test

// It was adapted from logxi package tests.

import (
	"bytes"
	"testing"
	"time"

	"github.com/grafov/kiwi"
	"github.com/grafov/kiwi/level"
	"github.com/grafov/kiwi/strict"
	"github.com/grafov/kiwi/timestamp"
)

// These tests write out all log levels with concurrency turned on and
// (mostly) equivalent fields.

func BenchmarkLevelsKiwiStrict_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.With(timestamp.Set(time.RFC3339))
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.AsJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Info(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Warn(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
		l.Error(strict.Int("key", 1), strict.Float64("key2", 3.141592), strict.String("key3", "string"), strict.Bool("key4", false))
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiStrictComplex_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.With(timestamp.Set(time.RFC3339))
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.AsJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Info(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Warn(strict.Int("key", 1), strict.Stringer("obj", testObject))
		l.Error(strict.Int("key", 1), strict.Stringer("obj", testObject))
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwi_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.With(timestamp.Set(time.RFC3339))
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.AsJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiComplex_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := level.New()
	l.With("_n", "bench", "_p", pid)
	l.With(timestamp.Set(time.RFC3339))
	level.LevelName = "l"
	out := kiwi.SinkTo(buf, kiwi.AsJSON()).Start()
	for i := 0; i < b.N; i++ {
		l.Debug("key", 1, "obj", testObject)
		l.Info("key", 1, "obj", testObject)
		l.Warn("key", 1, "obj", testObject)
		l.Error("key", 1, "obj", testObject)
	}
	b.StopTimer()
	out.Close()
}

func BenchmarkLevelsKiwiGlobal_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	out := kiwi.SinkTo(buf, kiwi.AsJSON()).Start()
	for i := 0; i < b.N; i++ {
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "debug", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "info", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "warn", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "error", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Close()
}
