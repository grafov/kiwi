package kiwi_test

// It was adapted from logxi package tests.

import (
	"bytes"
	"testing"
	"time"

	"github.com/grafov/kiwi"
)

// These tests write out all log levels with concurrency turned on and
// (mostly) equivalent fields.

func BenchmarkLevelsKiwiTyped_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwiTypedComplex_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwiTypedHelpers_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwiTypedHelpersComplex_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwi_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwiComplex_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
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
	out.Close()
}

func BenchmarkLevelsKiwiGlobal_JSON(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	out := kiwi.SinkTo(buf, kiwi.UseJSON()).Start()
	for i := 0; i < b.N; i++ {
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "debug", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "info", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "warn", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		kiwi.Log("t", time.Now().Format(time.RFC3339), "l", "error", "_n", "bench", "_p", pid, "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
	out.Close()
}
