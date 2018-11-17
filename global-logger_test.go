package kiwi_test

/*
Copyright (c) 2016-2018, Alexander I.Grafov <grafov@gmail.com>
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
	"github.com/grafov/kiwi"

	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Test logging of string value.
func TestGlobalLogger_LogStringValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", "The sample string with a lot of spaces.")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=\"The sample string with a lot of spaces.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of byte array.
func TestGlobalLogger_LogBytesValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", []byte("The sample string with a lot of spaces."))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=\"The sample string with a lot of spaces.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of integer value.
func TestGlobalLogger_LogIntValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", 123)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=123" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of negative integer value.
func TestGlobalLogger_LogNegativeIntValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", 123)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=123" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of float value in default (scientific) format.
func TestGlobalLogger_LogFloatValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", 3.14159265359)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=3.14159265359e+00" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of float value in fixed format.
func TestGlobalLogger_LogFixedFloatValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.FloatFormat = 'f'
	kiwi.Log("k", 3.14159265359)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=3.14159265359" {
		println(output.String())
		t.Fail()
	}
	// Turn back to default format.
	kiwi.FloatFormat = 'e'
}

// Test logging of boolean value.
func TestGlobalLogger_LogBoolValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", true, "k2", false)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=true k2=false" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of complex number.
func TestGlobalLogger_LogComplexValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("k", .12345E+5i, "k2", 1.e+0i)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "k=(0.000000+12345.000000i) k2=(0.000000+1.000000i)" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of time literal.
func TestGlobalLogger_LogTimeValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()
	value := time.Now()
	valueString := value.Format(kiwi.TimeLayout)

	kiwi.Log("k", value)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != fmt.Sprintf("k=%s", valueString) {
		println(output.String())
		t.Fail()
	}
}

// Test logging of the numeric key.
func TestGlobalLogger_LogNumericKey_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("123", "The sample value.")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "123=\"The sample value.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of the key with spaces.
func TestGlobalLogger_LogKeyWithSpaces_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("key with spaces", "The sample value.")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "\"key with spaces\"=\"The sample value.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of the key with tabs.
func TestGlobalLogger_LogKeyWithTabs_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log(fmt.Sprintf("key\twith\ttabs"), "The sample value.")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "\"key\\twith\\ttabs\"=\"The sample value.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of the multi lines key.
func TestGlobalLogger_LogKeyMultiLine_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log(fmt.Sprintf("multi\nlines\nkey"), "The sample value.")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "\"multi\\nlines\\nkey\"=\"The sample value.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test logging of the multi lines value.
func TestGlobalLogger_LogValueMultiLine_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.Log("key", fmt.Sprintf("multi\nlines\nvalue"))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != "key=\"multi\\nlines\\nvalue\"" {
		println(output.String())
		t.Fail()
	}
}

// Test log with the context value.
func TestGlobalLogger_WithContextPassed_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	kiwi.With("key1", "value")
	kiwi.Log("key2", "value")

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `key1="value" key2="value"` {
		t.Fail()
	}
}

// Test log with adding then removing the context.
func TestGlobalLogger_WithoutContextPassed_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	// add the context
	kiwi.With("key1", "value")
	// remove the context
	kiwi.Without("key1")
	// add regular pair
	kiwi.Log("key2", "value")

	out.Flush()
	if strings.TrimSpace(output.String()) != `key2="value"` {
		t.Fail()
	}
}

// Test log with adding then reset the context.
func TestGlobalLogger_ResetContext_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.AsLogfmt()).Start()

	// add the context
	kiwi.With("key1", "value")
	// reset the context
	kiwi.ResetContext()
	// add regular pair
	kiwi.Log("key2", "value")

	out.Flush()
	if strings.TrimSpace(output.String()) != `key2="value"` {
		t.Fail()
	}
}
