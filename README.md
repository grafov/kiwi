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
    BenchmarkLog-4                 20000         66733 ns/op       0.03 MB/s        7407 B/op        124 allocs/op
    BenchmarkLogComplex-4          20000         92365 ns/op       0.02 MB/s       11871 B/op        200 allocs/op
    BenchmarkKiwi-4                20000         82286 ns/op       0.02 MB/s        7410 B/op        130 allocs/op
    BenchmarkKiwiComplex-4         10000        150488 ns/op       0.01 MB/s       11176 B/op        169 allocs/op
    BenchmarkLogxi-4               50000         36533 ns/op       0.05 MB/s        4117 B/op         74 allocs/op
    BenchmarkLogxiComplex-4        20000         86944 ns/op       0.02 MB/s        9227 B/op        162 allocs/op
    BenchmarkLogrus-4              20000         97535 ns/op       0.02 MB/s       12168 B/op        184 allocs/op
    BenchmarkLogrusComplex-4       10000        114248 ns/op       0.02 MB/s       13943 B/op        238 allocs/op
    BenchmarkLog15-4               10000        139884 ns/op       0.01 MB/s       15618 B/op        224 allocs/op
    BenchmarkLog15Complex-4        10000        157494 ns/op       0.01 MB/s       14995 B/op        245 allocs/op

## Origins

* logfmt description [brandur.org/logfmt](https://brandur.org/logfmt)
* logfmt realization in Go and specs [godoc.org/github.com/kr/logfmt](https://godoc.org/github.com/kr/logfmt)
* picture used for logo [openclipart.org/detail/4416/kiwi-bird](https://openclipart.org/detail/4416/kiwi-bird)
