package where

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

All tests consists of three parts:

- arrange structures and initialize objects for use in tests
- act on testing object
- check and assert on results

These parts separated by empty lines in each test function.
*/

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grafov/kiwi"
)

func TestWhere_GetAllInfo_Logfmt(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := kiwi.New()
	out := kiwi.SinkTo(stream, kiwi.AsLogfmt()).Start()

	log.With(What(FilePos | Function))
	log.Log("key", "value")

	out.Flush().Close()
	expected := `file="`
	if !strings.Contains(stream.String(), expected) {
		t.Logf("expected %s got %v", expected, stream.String())
		t.Fail()
	}
	expected = `func="`
	if !strings.Contains(stream.String(), expected) {
		t.Logf("expected %s got %v", expected, stream.String())
		t.Fail()
	}
}

// Test of log to the stopped sink.
func TestWhereGlobal_GetAllInfo_Logfmt(t *testing.T) {
	stream := bytes.NewBufferString("")
	out := kiwi.SinkTo(stream, kiwi.AsLogfmt()).Start()

	kiwi.With(What(FilePos | Function))
	kiwi.Log("k", "v", "k2", "v2")

	out.Flush().Close()
	expected := `file="`
	if !strings.Contains(stream.String(), expected) {
		t.Logf("expected '%s' in '%v'", expected, stream.String())
		t.Fail()
	}
	expected = `function="`
	if !strings.Contains(stream.String(), expected) {
		t.Logf("expected '%s' in '%v'", expected, stream.String())
		t.Fail()
	}
}
