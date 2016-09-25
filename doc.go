// Copyright (c) 2016, Alexander I.Grafov aka Axel
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of kvlog nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// ॐ तारे तुत्तारे तुरे स्व

/*
Package kiwi is a library with an odd logic that log your application' data in its own strange way.

WIP. API and features is subject of changes. Use it carefully!

Features briefly:

 * simple format with explicit key for each log message (*logfmt* like) - for high readability by humans

 * optional JSON format that so liked by machines

 * there are not nailed levels, not hardcoded fields in the format

 * output dynamic filtering (change log verbosity on the fly)

 * can keep context of application

Key feature of `kiwi` logger is dynamic filtering of incoming records.
Instead of checking severety level for decide about pass or not the record to the output,
`kiwi` passes all records to *all* the outputs (they called *sinks* in `kiwi` terminology).
But before actual writing each record checked with a set of filters.
Each sink has its own filter set.
It takes into account record keys, values, ranges of values.
So each sink decides pass the record to a writer or filter it out.
Also any pairs in the record may be hidden: so different sinks may display different parts of the same record.
Other effect is: any record may be written to any number of outputs.
*/
package kiwi
