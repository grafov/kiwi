package kiwi

/*
Copyright (c) 2016-2019, Alexander I.Grafov <grafov@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of kvlog nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"testing"
)

/* All tests consists of three parts:

- arrange structures and initialize objects for use in tests
- act on testing object
- check and assert on results

These parts separated by empty lines in each test function.
*/

// Get context from logger. Helper for testing.
func (l *Logger) getAllContext() []*Pair {
	return l.context
}

func (l *Logger) checkContext(key string) string {
	for _, v := range l.context {
		if v.Key == key {
			return v.Val
		}
	}
	return ""
}

func TestNewLogger(t *testing.T) {
	l := New()

	if l == nil {
		t.Fatal("initialized logger is nil")
	}
}

// Test of creating a fork the logger from the existing logger.
func TestLogger_ForkWithContext(t *testing.T) {
	log := Fork().With("key", "value", "key2", 123)

	sub := log.Fork()

	if sub.checkContext("key") != "value" {
		t.Logf("expected %s got %v", "value", log.checkContext("key"))
		t.Fail()
	}
	if sub.checkContext("key2") != "123" {
		t.Logf("expected %s got %v", "123", log.checkContext("key2"))
		t.Fail()
	}
}

// Test of a new sublogger without context inheritance (because it
// should use Fork() for copy the context).
func TestLogger_NewWithContext(t *testing.T) {
	log := New().With("key", "value", "key2", 123)

	sub := log.New()

	context := sub.getAllContext()
	if len(context) != 0 {
		t.Logf("expected empty context but got %v", context)
		t.FailNow()
	}
}

// Test of creating a new logger from existing logger. Deleted values should not present in sublogger.
func TestLogger_ForkWithPartialContext(t *testing.T) {
	log := Fork().With("key", "value", "key2", "value2")

	log.Without("key")
	sub := log.Fork()

	if sub.getAllContext() == nil {
		t.Logf("context is nil but should not")
		t.FailNow()
	}
	if sub.checkContext("key") != "" {
		t.Logf(`expected nothing got %v`, sub.checkContext("key"))
		t.Fail()

	}
	if sub.checkContext("key2") == "" || sub.checkContext("key2") != "value2" {
		t.Logf(`expected "value2" got %v`, sub.checkContext("key2"))
		t.Fail()
	}
}
