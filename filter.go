package kiwi

// This file consists of implementations of Filter interface.

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
	"strconv"
)

type filter interface {
	Check(string, string) bool
}

type keyFilter struct {
}

func (*keyFilter) Check(key, val string) bool {
	return true
}

type valsFilter struct {
	Vals []string
}

func (k *valsFilter) Check(key, val string) bool {
	for _, v := range k.Vals {
		if v == val {
			return true
		}
	}
	return false
}

type rangeInt64Filter struct {
	From, To int64
}

func (k *rangeInt64Filter) Check(key, val string) bool {
	var (
		intVal int64
		err    error
	)
	if intVal, err = strconv.ParseInt(val, 10, 64); err != nil {
		return false
	}
	return intVal > k.From && intVal <= k.To
}

type rangeFloat64Filter struct {
	From, To float64
}

func (k *rangeFloat64Filter) Check(key, val string) bool {
	var (
		floatVal float64
		err      error
	)
	if floatVal, err = strconv.ParseFloat(val, 64); err != nil {
		return false
	}
	return floatVal > k.From && floatVal <= k.To
}
