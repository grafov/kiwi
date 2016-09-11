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

// FormatLogfmt formats filtered output as Logfmt format.
func FormatLogfmt(record []pair) []byte {
	var line bytes.Buffer
	for _, pair := range record {
		line.WriteString(pair.Key)
		line.WriteRune('=')
		if pair.Val.Quoted {
			line.WriteString(strconv.Quote(pair.Val.Strv))
		} else {
			line.WriteString(pair.Val.Strv)
		}
		line.WriteRune(' ')
	}
	return line.Bytes()
}

// FormatJSON formats filtered output as JSON.
// Function accepts slice of record pairs of Pair type.
// The function shall not modify record pairs.
func FormatJSON(record []pair) []byte {
	var line bytes.Buffer
	line.WriteRune('{')
	for _, pair := range record {
		line.WriteRune('"')
		line.WriteString(pair.Key)
		line.WriteString("\":")
		if pair.Val.Quoted {
			line.WriteString(strconv.Quote(pair.Val.Strv))
		} else {
			line.WriteString(pair.Val.Strv)
		}
		line.WriteString(", ")
	}
	line.WriteRune('}')
	return line.Bytes()
}
