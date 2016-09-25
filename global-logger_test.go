package kiwi_test

/*
Copyright (c) 2016, Alexander I.Grafov aka Axel
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
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", "The sample string with a lot of spaces.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=\"The sample string with a lot of spaces.\"" {
		t.Fail()
	}
}

// Test logging of byte array.
func TestGlobalLogger_LogBytesValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", []byte("The sample string with a lot of spaces."))

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=\"The sample string with a lot of spaces.\"" {
		t.Fail()
	}
}

// Test logging of integer value.
func TestGlobalLogger_LogIntValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", 123)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=123" {
		t.Fail()
	}
}

// Test logging of negative integer value.
func TestGlobalLogger_LogNegativeIntValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", 123)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=123" {
		t.Fail()
	}
}

// Test logging of float value in default (scientific) format.
func TestGlobalLogger_LogFloatValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", 3.14159265359)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=3.14159265359e+00" {
		t.Fail()
	}
}

// Test logging of float value in fixed format.
func TestGlobalLogger_LogFixedFloatValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.FloatFormat = 'f'
	kiwi.Log("k", 3.14159265359)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=3.14159265359" {
		t.Fail()
	}
	// Turn back to default format.
	kiwi.FloatFormat = 'e'
}

// Test logging of boolean value.
func TestGlobalLogger_LogBoolValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", true, "k2", false)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=true k2=false" {
		t.Fail()
	}
}

// Test logging of complex number.
func TestGlobalLogger_LogComplexValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	kiwi.Log("k", .12345E+5i, "k2", 1.e+0i)

	out.Flush()
	if strings.TrimSpace(output.String()) != "k=(0.000000+12345.000000i) k2=(0.000000+1.000000i)" {
		t.Fail()
	}
}

// Test logging of time literal.
func TestGlobalLogger_LogTimeValue_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	value := time.Now()
	valueString := value.Format(kiwi.TimeLayout)
	defer out.Close()

	kiwi.Log("k", value)

	out.Flush()
	if strings.TrimSpace(output.String()) != fmt.Sprintf("k=%s", valueString) {
		t.Fail()
	}
}
