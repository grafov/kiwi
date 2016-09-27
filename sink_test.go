package kiwi

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
	"bytes"
	"strings"
	"testing"
)

/* All tests consists of three parts:

- arrange structures and initialize objects for use in tests
- act on testing object
- check and assert on results

These parts separated by empty lines in each test function.
*/

// Test of log to the stopped sink.
func TestSink_LogToStoppedSink_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt())
	defer out.Close()

	log.Log("k", "The sample string that should be ignored.")

	out.Flush()
	if strings.TrimSpace(output.String()) != "" {
		t.Fail()
	}
}

// Test of log to the stopped sink. It should not crash logger.
func TestSink_StopTwice(t *testing.T) {
	out := SinkTo(bytes.NewBufferString(""), UseLogfmt())
	out.Stop()
	out.Close()
}

// Test of the close already closed sink. It should not crash logger.
func TestSink_CloseTwice(t *testing.T) {
	out := SinkTo(bytes.NewBufferString(""), UseLogfmt())
	out.Close()
	out.Close()
}

// Test of WithKey filter. It should pass record to the output.
func TestSink_WithKeyFilterPass(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).WithKey("Gandalf").Start()
	defer out.Close()

	log.Log("Gandalf", "You shall not pass!") // from the movie

	out.Flush()
	if strings.TrimSpace(output.String()) != "Gandalf=\"You shall not pass!\"" {
		t.Fail()
	}

}

// Test of WithoutKey filter. It should not pass record to the output.
func TestSink_WithoutKeyFilterOut(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).WithoutKey("Gandalf").Start()
	defer out.Close()

	log.Log("Gandalf", "You cannot pass!") // from the book

	out.Flush()
	if strings.TrimSpace(output.String()) != "" {
		t.Fail()
	}
}

// Test of WithValue filter. It should pass the record to the output because the key missed.
func TestSink_WithValueFilterMissedKeyPass(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).WithValue("Gandalf", "You cannot pass!").Start()
	defer out.Close()

	log.Log("Balrog", "Boo!")

	out.Flush()
	if strings.TrimSpace(output.String()) != "Balrog=\"Boo!\"" {
		t.Fail()
	}

}
