package kiwi

/*
Copyright (c) 2016, Alexander I.Grafov <grafov@gmail.com>
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
	"time"
)

// Test of log to the stopped sink.
func TestSink_LogToStoppedSink_Logfmt(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt())

	log.Log("key", "The sample string that should be ignored.")

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of log to the stopped sink.
func TestSink_LogToStoppedSink_JSON(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseJSON())

	log.Log("key", "The sample string that should be ignored.")

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of log to the stopped sink. It should not crash logger.
func TestSink_StopTwice(t *testing.T) {
	out := SinkTo(bytes.NewBufferString(""), UseLogfmt())
	out.Stop()
	out.Close()
}

// Test of start already started sink. It should not crash logger.
func TestSink_StartTwice(t *testing.T) {
	out := SinkTo(bytes.NewBufferString(""), UseLogfmt()).Start()
	out.Start()
	out.Close()
}

// Test of the close already closed sink. It should not crash logger.
func TestSink_CloseTwice(t *testing.T) {
	out := SinkTo(bytes.NewBufferString(""), UseLogfmt())
	out.Close()
	out.Close()
}

// Test of reuse of the already created sink.
func TestSink_SinkReuse(t *testing.T) {
	stream := bytes.NewBufferString("")
	out := SinkTo(stream, UseLogfmt())

	SinkTo(stream, UseJSON())
	SinkTo(stream, UseLogfmt())

	out.Close()
}

// Test of HideKey. It should pass record to the output.
func TestSink_HideKey(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt())

	out.Start().Hide("two")
	log.Log("one", 1, "two", 2, "three", 3)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `one=1 three=3` {
		println(stream.String())
		t.Fail()
	}
}

// Test of UnhideKey. It should pass record to the output.
func TestSink_UnhideKey(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt())

	out.Hide("two").Start().Unhide("two")
	log.Log("one", 1, "two", 2, "three", 3)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `one=1 two=2 three=3` {
		println(stream.String())
		t.Fail()
	}
}

// Test of unhide already unhidden key. It should pass record to the output.
func TestSink_UnhideKeyTwice(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt())

	out.Start().Unhide("one").Unhide("one")
	log.Log("one", 1, "two", 2)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `one=1 two=2` {
		println(stream.String())
		t.Fail()
	}

}

// Test of WithKey filter. It should pass record to the output.
func TestSink_WithKeyFilterPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithKey("Gandalf").Start()

	log.Log("Gandalf", "You shall not pass!") // cite from the movie

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `Gandalf="You shall not pass!"` {
		println(stream.String())
		t.Fail()
	}

}

// Test of WithoutKey filter. It should not pass record to the output.
func TestSink_WithoutKeyFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithoutKey("Gandalf").Start()

	log.Log("Gandalf", "You cannot pass!") // cite from the book

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithValue filter. It should pass the record to the output because the key missed.
func TestSink_WithValueFilterMissedKeyPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithValue("key", "passed").Start()

	log.Log("key", "passed")

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key="passed"` {
		println(stream.String())
		t.Fail()
	}

}

// Test of WithValue filter. It should pass the record to the output because the value matched.
func TestSink_WithValueFilterPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithValue("key", "passed", "and this passed too").Start()

	log.Log("key", "passed", "key", "and this passed too")

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key="passed" key="and this passed too"` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithValue filter. It should filter out the record because no one value matched.
func TestSink_WithValueFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithValue("key", "filtered", "out").Start()

	log.Log("key", "try it")

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithIntRange filter. It should pass the record to the output because the key missed.
func TestSink_WithIntRangeFilterMissedKeyPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithInt64Range("key", 1, 2).Start()

	log.Log("another key", 3)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `"another key"=3` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithIntRange filter. It should pass the record to the output because the value in the range.
func TestSink_WithIntRangeFilterPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithInt64Range("key", 1, 3).Start()

	log.Log("key", 2)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key=2` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithIntRange filter. It should filter out the record because the value not in the range.
func TestSink_WithIntRangeFilterFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithInt64Range("key", 1, 3).Start()

	log.Log("key", 4)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithFloatRange filter. It should pass the record to the output because the key missed.
func TestSink_WithFloatRangeFilterMissedKeyPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithFloat64Range("key", 1.0, 2.0).Start()

	log.Log("another key", 3)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `"another key"=3` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithFloatRange filter. It should pass the record to the output because the value in the range.
func TestSink_WithFloatRangeFilterPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithFloat64Range("key", 1.0, 3.0).Start()

	log.Log("key", 2.0)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key=2e+00` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithFloatRange filter. It should filter out the record because the value not in the range.
func TestSink_WithFloatRangeFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	out := SinkTo(stream, UseLogfmt()).WithFloat64Range("key", 1.0, 3.0).Start()

	log.Log("key", 4.0)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithTimeRange filter. It should pass the record to the output because the value in the range.
func TestSink_WithTimeRangeFilterPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	now := time.Now()
	hourAfterNow := now.Add(1 * time.Hour)
	halfHourAfterNow := now.Add(30 * time.Minute)
	halfHourAsString := halfHourAfterNow.Format(TimeLayout)
	out := SinkTo(stream, UseLogfmt()).WithTimeRange("key", now, hourAfterNow).Start()

	log.Log("key", halfHourAfterNow)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key=`+halfHourAsString {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithTimeRange filter. It should filter out the record because the value not in the range.
func TestSink_WithTimeRangeFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	now := time.Now()
	hourAfterNow := now.Add(1 * time.Hour)
	halfHourAfterNow := now.Add(30 * time.Minute)
	out := SinkTo(stream, UseLogfmt()).WithTimeRange("key", now, halfHourAfterNow).Start()

	log.Log("key", hourAfterNow)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}

type customFilterThatReturnsTrue struct{}

func (customFilterThatReturnsTrue) Check(key, val string) bool {
	return true
}

// Test of WithFilter custom filter. It should pass the record to the output because the value in the range.
func TestSink_WithCustomPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	var customFilter customFilterThatReturnsTrue
	out := SinkTo(stream, UseLogfmt()).WithFilter("key", customFilter).Start()

	log.Log("key", 2)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `key=2` {
		println(stream.String())
		t.Fail()
	}
}

// Test of WithFilter custom filter. It should pass the record to the output because the key missed.
func TestSink_WithCustomMissedKeyPass(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	var customFilter customFilterThatReturnsTrue
	out := SinkTo(stream, UseLogfmt()).WithFilter("key", customFilter).Start()

	log.Log("another key", 3)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != `"another key"=3` {
		println(stream.String())
		t.Fail()
	}
}

type customFilterThatReturnsFalse struct{}

func (customFilterThatReturnsFalse) Check(key, val string) bool {
	return false
}

// Test of WithFilter custom filter. It should pass the record to the output because the value in the range.
func TestSink_WithCustomFilterOut(t *testing.T) {
	stream := bytes.NewBufferString("")
	log := New()
	var customFilter customFilterThatReturnsFalse
	out := SinkTo(stream, UseLogfmt()).WithFilter("key", customFilter).Start()

	log.Log("key", 2)

	out.Flush().Close()
	if strings.TrimSpace(stream.String()) != "" {
		println(stream.String())
		t.Fail()
	}
}
