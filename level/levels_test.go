package level

/* Copyright (c) 2016, Alexander I.Grafov aka Axel
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

// Test of log with fatal level with empty value. Useless but function allow it.
func TestLoggerLevels_LogFatalEmpty_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Fatal()

	out.Flush()
	if strings.TrimSpace(output.String()) != "level=\"fatal\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with fatal level without a key.
func TestLoggerLevels_LogFatal_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Fatal("The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "level=\"fatal\" \"The sample message.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with fatal level with a key.
func TestLoggerLevels_LogFatalWKey_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Fatal("msg", "The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "msg=\"The sample message.\" level=\"fatal\"" {
		t.Fail()
	}
}

// Test of log with critical level without a key.
func TestLoggerLevels_LogCrit_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Crit("The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "level=\"critical\" \"The sample message.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with critical level with a key.
func TestLoggerLevels_LogCritWKey_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Crit("msg", "The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "msg=\"The sample message.\" level=\"critical\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with error level without a key.
func TestLoggerLevels_LogError_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Error("The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "level=\"error\" \"The sample message.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with error level with a key.
func TestLoggerLevels_LogErrorWKey_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Error("msg", "The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "msg=\"The sample message.\" level=\"error\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with warning level without a key.
func TestLoggerLevels_LogWarn_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Warn("The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "level=\"warning\" \"The sample message.\"" {
		println(output.String())
		t.Fail()
	}
}

// Test of log with warning level with a key.
func TestLoggerLevels_LogWarnWKey_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := kiwi.SinkTo(output, kiwi.UseLogfmt()).Start()
	defer out.Close()

	log.Warn("msg", "The sample message.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "msg=\"The sample message.\" level=\"warning\"" {
		println(output.String())
		t.Fail()
	}
}
