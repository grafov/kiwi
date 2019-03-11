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
	"bytes"
	"sync"
	"testing"
)

/* All tests consists of three parts:

- arrange structures and initialize objects for use in tests
- act on testing object
- check and assert on results

These parts separated by empty lines in each test function.
*/

// Test logging of string value.
func TestRace_NewFromGlobal_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := SinkTo(output, AsLogfmt()).Start()
	defer out.Close()

	for i := 0; i < 100; i++ {
		go func(i int) {
			log := New()
			log.Log("k", "The sample string.", "instance", i)
		}(i)
	}
	Log("k", "The sample string.", "instance", 0)

	out.Flush()
}

// Test logging of string value.
func TestRace_ForkFromGlobal_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := SinkTo(output, AsLogfmt()).Start()
	defer out.Close()

	for i := 0; i < 100; i++ {
		go func(i int) {
			log := Fork()
			log.Log("k", "The sample string.", "instance", i)
		}(i)
	}
	Log("k", "The sample string.", "instance", 0)

	out.Flush()
}

// Test logging of string value.
func TestRace_Fork_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := SinkTo(output, AsLogfmt()).Start()
	defer out.Close()
	var (
		wg sync.WaitGroup
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			log := Fork()
			log.Log("k", "The sample string.", "instance", i)
			wg.Done()
		}(i)
	}
	Log("k", "The sample string.", "instance", 0)
	wg.Wait()

	out.Flush()
}

// Test logging of string value.
func TestRace_ForkFork_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := SinkTo(output, AsLogfmt()).Start()
	defer out.Close()
	var (
		wg sync.WaitGroup
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			log := Fork().Fork()
			log.Log("k", "The sample string.", "instance", i)
			wg.Done()
		}(i)
	}
	Log("k", "The sample string.", "instance", 0)
	wg.Wait()

	out.Flush()
}

// Test logging of string value.
func TestRace_ForkInside_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	out := SinkTo(output, AsLogfmt()).Start()
	defer out.Close()
	var (
		wg  sync.WaitGroup
		log = Fork()
	)
	log.With("logger", "global")
	log.Log("k", "The sample string.", "instance", 0)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int, log *Logger) {
			l := log.Fork().With("logger", i)
			l.Log("k", "The sample string.", "instance", i)
			wg.Done()
		}(i, log)
	}
	log.Log("k", "The sample string.", "instance", 0)
	wg.Wait()

	out.Flush()
}

// // Test logging of byte array.
// func TestLogger_LogBytesValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", []byte("The sample string with a lot of spaces."))

//	out.Flush()
//	if strings.TrimSpace(output.String()) != `k="The sample string with a lot of spaces."` {
//		t.Fail()
//	}
// }

// // Test logging of integer value.
// func TestLogger_LogIntValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", 123)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=123" {
//		t.Fail()
//	}
// }

// // Test logging of negative integer value.
// func TestLogger_LogNegativeIntValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", -123)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=-123" {
//		t.Fail()
//	}
// }

// // Test logging of float value in default (scientific) format.
// func TestLogger_LogFloatValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", 3.14159265359)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=3.14159265359e+00" {
//		t.Fail()
//	}
// }

// // Test logging of float value in fixed format.
// func TestLogger_LogFixedFloatValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	FloatFormat = 'f'
//	log.Log("k", 3.14159265359)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=3.14159265359" {
//		t.Fail()
//	}
//	// Turn back to default format.
//	FloatFormat = 'e'
// }

// // Test logging of boolean value.
// func TestLogger_LogBoolValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", true, "k2", false)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=true k2=false" {
//		t.Fail()
//	}
// }

// // Test logging of complex number.
// func TestLogger_LogComplexValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log("k", .12345E+5i, "k2", 1.e+0i)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != "k=(0.000000+12345.000000i) k2=(0.000000+1.000000i)" {
//		t.Fail()
//	}
// }

// // Test logging of time literal.
// func TestLogger_LogTimeValue_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	value := time.Now()
//	valueString := value.Format(TimeLayout)
//	defer out.Close()

//	log.Log("k", value)

//	out.Flush()
//	if strings.TrimSpace(output.String()) != fmt.Sprintf("k=%s", valueString) {
//		t.Fail()
//	}
// }

// func TestLogger_LogIntKeyInvalid_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log(123, 456)

//	out.Flush()
//	expected := `kiwi-error="non a string type (int) for the key (123)" message=456`
//	if strings.TrimSpace(output.String()) != expected {
//		t.Logf("expected %s got %v", expected, output.String())
//		t.Fail()
//	}
// }

// // The logger accepts string keys. In other cases it complains about
// // the wrong key type. It also assumes that each key followed by the
// // value.
// func TestLogger_LogThreeIntKeysInvalid_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log(123, 456, 789)

//	out.Flush()
//	expected := `kiwi-error="non a string type (int) for the key (123)" message=456 kiwi-error="non a string type (int) for the key (789)"`
//	if strings.TrimSpace(output.String()) != expected {
//		t.Logf("expected %s got %v", expected, output.String())
//		t.Fail()
//	}
// }

// // The logger accepts string keys. In other cases it complains about
// // the wrong key type. It also assumes that each key followed by the
// // value.
// func TestLogger_LogFourIntKeysInvalid_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Log(12, 34, 56, 78)

//	out.Flush()
//	expected := `kiwi-error="non a string type (int) for the key (12)" message=34 kiwi-error="non a string type (int) for the key (56)" message=78`
//	if strings.TrimSpace(output.String()) != expected {
//		t.Logf("expected %s got %v", expected, output.String())
//		t.Fail()
//	}
// }

// // Test chaining for Add()
// func TestLogger_AddMixChained_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.Add("k", "value2").Add("k2", 123).Add("k3", 3.14159265359).Log()

//	out.Flush()
//	expected := `k="value2" k2=123 k3=3.14159265359e+00`
//	if strings.TrimSpace(output.String()) != expected {
//		t.Logf("expected %s got %v", expected, output.String())
//		t.Fail()
//	}
// }

// // Test log with the context value.
// func TestLogger_WithContextPassed_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	log.With("key1", "value")
//	log.Log("key2", "value")

//	out.Flush()
//	if strings.TrimSpace(output.String()) != `key1="value" key2="value"` {
//		t.Fail()
//	}
// }

// // Test log with adding then removing the context.
// func TestLogger_WithoutContextPassed_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	// add the context
//	log.With("key1", "value")
//	// add regular pair
//	log.Add("key2", "value")
//	// remove the context and flush the record
//	log.Without("key1").Log()

//	out.Flush()
//	if strings.TrimSpace(output.String()) != `key2="value"` {
//		t.Fail()
//	}
// }

// // Test log with adding then reset the context.
// func TestLogger_ResetContext_Logfmt(t *testing.T) {
//	output := bytes.NewBufferString("")
//	log := New()
//	out := SinkTo(output, AsLogfmt()).Start()
//	defer out.Close()

//	// add the context
//	log.With("key1", "value")
//	// add regular pair
//	log.Add("key2", "value")
//	// reset the context and flush the record
//	log.ResetContext().Log()

//	out.Flush()
//	if strings.TrimSpace(output.String()) != `key2="value"` {
//		t.Fail()
//	}
// }
