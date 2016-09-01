# Kiwi logger & context keeper [![Go Report Card](https://goreportcard.com/badge/grafov/kiwi)](https://goreportcard.com/report/grafov/kiwi)

*WIP. It works but not all features completed.*

![Kiwi bird](flomar-kiwi-bird-300px.png)

*Kiwi* /kiːwiː/ are birds native to New Zealand, in the genus Apteryx and family Apterygidae. They are flightless, have hair-like feathers and smell like mushrooms. They look strange and funny so when I wrote a logger for Go language I decided to devote it to this beast which I never seen in a wild (I live very far from places where kiwis are live).

*Kiwi Logger* — this is a library with an odd logic that log your application data in its own strange way.

## Features

Shortly: both humans and robots are love it!

Briefly:

* structured logging for high readability by humans
* optional JSON format that liked by machines
* dynamically selectable outputs (changing log verbosity on the fly)
* no nailed levels, no hardcoded fields in the format
* can keep context of application

## Docs [![GoDoc](https://godoc.org/github.com/grafov/kiwi?status.svg)](https://godoc.org/github.com/grafov/kiwi)

Examples of logger usage see at [cmd/*](cmd) subfolders.
See API description and code samples in [godoc](http://godoc.org/github.com/grafov/kiwi).

## Installation

Package have not external dependencies except standard library. So just

    go get github.com/grafov/kiwi

## Evaluating rules of record values

* Logged values evaluated *immediately* when they added to a record.
* Context values evaluated *once* when they added to a logger.
* For lazy evaluating of context and record values use workaround with functions without call them in a log record:

        # For lazy evaluating you need function that returns interface{} or []interface{}
        func longActionForDelayedEvaluation() interface{} {
           // do something complex
           return "something"
        }
        myLog.Add("lazy-sample", longActionForDelayedEvaluation) # but not longActionForDelayedEvaluation()

Logger recognizes next function types when adding key-val pairs to a record:

* `func () string`
* `func () interface{}`
* `func () []interface{}`

Hence value of `lazy-sample` from the example above will be evaluated only on `Log()` call.


## Instead of FAQ

0. Kiwi logger not strictly follow logfmt specs.
1. Ideas of key-value format very near to JSON output but with orientation on readability for humans without additional tools for log parsing.
2. Yes, it was architectured and developed to be a standard number 15 that competing with others. It is not pretend to be log format for everything.

## Similar works

* [log15](https://github.com/inconshreveable/log15)

## Comparison with other loggers

    $ go test -bench=. -benchmem
    BenchmarkLevelsKiwi-4                  50000         31437 ns/op       0.06 MB/s        6993 B/op        115 allocs/op
    BenchmarkLevelsKiwiComplex-4           30000         57841 ns/op       0.03 MB/s       10303 B/op        173 allocs/op
    BenchmarkLevelsStdLog-4               100000         21769 ns/op       0.09 MB/s        7159 B/op        124 allocs/op
    BenchmarkLevelsStdLogComplex-4         50000         35168 ns/op       0.06 MB/s       11446 B/op        200 allocs/op
    BenchmarkLevelsLogxi-4                100000         13628 ns/op       0.15 MB/s        4127 B/op         74 allocs/op
    BenchmarkLevelsLogxiComplex-4          50000         31377 ns/op       0.06 MB/s        8713 B/op        162 allocs/op
    BenchmarkLevelsLogrus-4                50000         38582 ns/op       0.05 MB/s       12320 B/op        177 allocs/op
    BenchmarkLevelsLogrusComplex-4         30000         45749 ns/op       0.04 MB/s       13990 B/op        231 allocs/op
    BenchmarkLevelsLog15-4                 30000         53837 ns/op       0.04 MB/s       14998 B/op        224 allocs/op
    BenchmarkLevelsLog15Complex-4          20000         61730 ns/op       0.03 MB/s       15127 B/op        246 allocs/op

## Origins

* logfmt description [brandur.org/logfmt](https://brandur.org/logfmt)
* logfmt realization in Go and specs [godoc.org/github.com/kr/logfmt](https://godoc.org/github.com/kr/logfmt)
* picture used for logo [openclipart.org/detail/4416/kiwi-bird](https://openclipart.org/detail/4416/kiwi-bird)
