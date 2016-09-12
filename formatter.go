package kiwi

// This file consists of realizations of default formatters.

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

ॐ तारे तुत्तारे तुरे स्व */

import (
	"bytes"
	"strconv"
)

// Formatter represents format of the output.
type Formatter interface {
	Begin()
	Pair(key, val string, quoted bool)
	Finish() []byte
}

type formatLogfmt struct {
	line bytes.Buffer
}

func UseLogfmt() *formatLogfmt {
	return new(formatLogfmt)
}

func (f *formatLogfmt) Begin() {
	f.line.Reset()
}

func (f *formatLogfmt) Pair(key, val string, quoted bool) {
	f.line.WriteString(key)
	f.line.WriteRune('=')
	if quoted {
		f.line.WriteString(strconv.Quote(val))
	} else {
		f.line.WriteString(val)
	}
	f.line.WriteRune(' ')
}

func (f *formatLogfmt) Finish() []byte {
	return f.line.Bytes()
}

type formatJSON struct {
	line bytes.Buffer
}

func UseJSON() *formatJSON {
	return new(formatJSON)
}

func (f *formatJSON) Begin() {
	f.line.Reset()
	f.line.WriteRune('{')
}

func (f *formatJSON) Pair(key, val string, quoted bool) {
	f.line.WriteRune('"')
	f.line.WriteString(key)
	f.line.WriteString("\":")
	if quoted {
		f.line.WriteString(strconv.Quote(val))
	} else {
		f.line.WriteString(val)
	}
	f.line.WriteString(", ")
}

func (f *formatJSON) Finish() []byte {
	f.line.WriteRune('}')
	return f.line.Bytes()
}
