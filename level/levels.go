package level

// This file consists of Logger methods for imitating oldschool logging with levels.

/* Copyright (c) 2016-2017, Alexander I.Grafov
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

import "github.com/grafov/kiwi"

type Logger struct {
	*kiwi.Logger
}

// New creates a new leveled logger instance.
func New() *Logger {
	return &Logger{kiwi.New()}
}

// Fatal imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "fatal". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any value what you want.
func (l *Logger) Fatal(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "fatal", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "fatal")...)
	}
}

// Crit imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "critical". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any value what you want.
func (l *Logger) Crit(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "critical", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "critical")...)
	}
}

// Error imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "error". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Error(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "error", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "error")...)
	}
}

// Warn imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "warning". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any recVal you want.
func (l *Logger) Warn(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "warning", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "warning")...)
	}
}

// Info imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "info". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any value what you want.
func (l *Logger) Info(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "info", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "info")...)
	}
}

// Debug imitates behaviour of common loggers with severity levels. It adds a record
// with severity "level" = "debug". Default severity name "level" may be changed
// globally for all package with UseLevelName(). There is nothing special in "level"
// key so it may be overrided with any value what you want.
func (l *Logger) Debug(keyVals ...interface{}) {
	if len(keyVals) == 1 {
		l.Log(LevelName, "debug", keyVals[0])
	} else {
		l.Log(append(keyVals, LevelName, "debug")...)
	}
}
