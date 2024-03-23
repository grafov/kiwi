package kiwi

// Obsoleted functions that will be removed just after v1.0.

/* Copyright (c) 2016-2024, Alexander I.Grafov <grafov@inet.name>
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

// FlushAll should wait for all the sinks to be flushed. It does
// nothing currently. It has left for compatibility with old API.
// Deprecated: just no need for this feature.
func FlushAll() {
	Log(InfoKey, "deprecated", MessageKey, "FlushAll() is deprecated, you should remove it from the code")
}

// Flush waits that all previously sent to the output records
// worked. It does nothing currently. It has left for compatibility
// with old API.
// Deprecated: just no need for sinks anymore.
func (s *Sink) Flush() *Sink {
	Log(InfoKey, "deprecated", MessageKey, "sink.Flush() is deprecated, you should remove it from the code")
	return s
}

// ResetContext resets the context of the logger. The function is not
// concurrent safe.
// Deprecated: use Reset() instead.
func (l *Logger) ResetContext() *Logger {
	Log(InfoKey, "deprecated", MessageKey, "ResetContext() is deprecated, you should remove it from the code")
	l.c.Lock()
	defer l.c.Unlock()
	l.context = nil
	return l
}

// ResetContext resets the global context for the global logger and
// its descendants. It is safe for concurrency.
// Deprecated: use Reset() instead.
func ResetContext() {
	Log(InfoKey, "deprecated", MessageKey, "ResetContext() is deprecated, you should remove it from the code")
	global.Lock()
	defer global.Unlock()
	context = nil
}

// WithoutAll resets the context of the logger. The function is not
// concurrent safe.
// Deprecated: use Reset() instead.
func (l *Logger) WithoutAll() *Logger {
	Log(InfoKey, "deprecated", MessageKey, "WithoutAll() is deprecated, you should remove it from the code")
	l.c.Lock()
	defer l.c.Unlock()
	l.context = nil
	return l
}

// WithoutAll resets the global context for the global logger and
// its descendants. It is safe for concurrency.
// Deprecated: use Reset() instead.
func WithoutAll() {
	Log(InfoKey, "deprecated", MessageKey, "WithoutAll() is deprecated, you should remove it from the code")
	global.Lock()
	context = nil
	global.Unlock()
}
