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

// Test of non default value of FloatFormat global var.
// It should format the value accordingly with selected format ('f' in this test).
func TestConvertor_NonDefaultFloatFormatPass_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	original := FloatFormat
	FloatFormat = 'f'
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log("key", 3.14159265)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `key=3.14159265` {
		println(output.String())
		t.Fail()
	}
	FloatFormat = original
}

// Test of non default value of TimeLayout global var.
// It should format the value accordingly with selected format.
func TestConvertor_NonDefaultTimeLayoutPass_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	original := TimeLayout
	TimeLayout = time.RFC822
	now := time.Now()
	nowString := now.Format(time.RFC822)
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log("key", now)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `key=`+nowString {
		println(output.String())
		t.Fail()
	}
	TimeLayout = original
}

func TestConvertor_LogByteType_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log([]byte("the key"), []byte("the sample byte sequence..."))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `"the key"="the sample byte sequence..."` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogBoolType_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(true, false)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `true=false` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogInt8Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(int8(1), int8(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogInt16Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(int16(1), int16(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogInt32Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(int32(1), int32(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogIntType_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(1, 2)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogInt64Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(int64(1), int64(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogUint8Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(uint8(1), uint8(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogUint16Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(uint16(1), uint16(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogUint32Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(uint32(1), uint32(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogUintType_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(uint(1), uint(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogUint64Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(uint64(1), uint64(2))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `1=2` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogFloat32Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(float32(3.14159265), float32(3.14159265))

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `3.1415927e+00=3.1415927e+00` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogFloat64Type_Logfmt(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(3.14159265, 3.14159265)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `3.14159265e+00=3.14159265e+00` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogNil(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(nil)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `message="<nil>"` {
		println(output.String())
		t.Fail()
	}
}

func TestConvertor_LogNilNil(t *testing.T) {
	output := bytes.NewBufferString("")
	log := New()
	out := SinkTo(output, UseLogfmt()).Start()

	log.Log(nil, nil)

	out.Flush().Close()
	if strings.TrimSpace(output.String()) != `<nil>="<nil>"` {
		println(output.String())
		t.Fail()
	}
}
