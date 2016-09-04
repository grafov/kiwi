# Kiwi logger & context keeper [![Go Report Card](https://goreportcard.com/badge/grafov/kiwi)](https://goreportcard.com/report/grafov/kiwi)

*WIP. It works but not all features completed.*

![Kiwi bird](flomar-kiwi-bird-300px.png)

*Kiwi* /kiːwiː/ are birds native to New Zealand, in the genus Apteryx and family Apterygidae. They are flightless, have hair-like feathers and smell like mushrooms. They look strange and funny so when I wrote a logger for Go language I decided to devote it to this beast which I never seen in a wild (I live very far from places where kiwis are live).

*Kiwi Logger* — this is a library with an odd logic that log your application data in its own strange way.

## Features

Shortly: both humans and robots will love it!

Briefly:

* structured logging for high readability by humans
* optional JSON format that liked by machines
* dynamically selectable outputs (changing log verbosity on the fly)
* no nailed levels, no hardcoded fields in the format
* can keep context of application

Kiwi logger has two primary design goals:

1. Convenient structured logging syntax: logfmt as default format, method chaining.
2. Separation of logging flow from control flow: you log everything without conditions but output filtering really select what and where will be saved.

Point 2 needs more explanations. 
Traditional way for logging is set level of severity for each log record. 
And check level before passing this record to a writer.
It is not bad way but it is not obvious. 
Especially when logger introduces many severity levels like "debug", "info", "warning", "critical", "fatal", "panic" and so on. 
Look the internet for many guides with controversial recommendations how to distinguish all these "standard" levels and apply them to various events in your application.
When you should use "fatal" instead of "panic" or "debug" instead of "info".

There is alternative way not use severity levels at all. Structured logging do in this way. 
Especially it right for *logfmt* format.
It is not required part of logfmt or structured logging but it naturally ensued from rules they offer.
You just log pairs of keys and values and these pairs may be of any kind. There is not a standard list of keys.
If you need you can use levels for log records for example set key named "level" (or any other name you want) with values INFO, WARNING, ERROR etc.
But it is not requirement.
So you log just pairs of arbitrary keys and values and interprete them as you wish.

Feature of kiwi logger is dynamic filtering of incoming records.
You log all then set for each output point what you want to see in log stream.
And it can be changed in any moment:
`kiwi` has methods for filtering by keys, values, ranges of values.
Use case for it is export some handler for setting these filters and you got
ability dynamycally change flow and verbosity of logs.
For example increase verbosity for a specific module or a handler and decrease for rest of the application.

## Docs [![GoDoc](https://godoc.org/github.com/grafov/kiwi?status.svg)](https://godoc.org/github.com/grafov/kiwi)

Examples of logger usage see at [cmd/*](cmd) subfolders.
See API description and code samples in [godoc](http://godoc.org/github.com/grafov/kiwi).

## Installation [![Build Status](https://travis-ci.org/grafov/kiwi.svg?branch=master)](https://travis-ci.org/grafov/kiwi)

Package have not external dependencies except standard library. So just

    go get github.com/grafov/kiwi

## Usage examples

```go
import "github.com/grafov/kiwi"

func main() {
	// Creates a new logger instance.
	log:=kiwi.New()

	// Now just log something as key/value pair. It will pass to output immediately (read about outputs below).
	log.Log("msg", "something", "another key", "another value")
	// Expected output:
	// msg="something" another\ key="another value"

	// You can pass odd number of parameters. Odd parameter passed to output just as is.
	log.Log("key-only")
	// Expected output:
	// "key-only"

	// It can add key=value pairs to a new log record.
	// They don't passed to the output until Log() call.
	log.Add("where", "module1", "event", "something happened")

	// So it may be any number of Add() calls with additional pairs.
	// Then flush them all.
	log.Add("event", "and now something completely different").Log()

	// You can pass any scalar types from Go standard library as record keys and values
	// they will be converted to their string representation.
	log.Log("key", 123, "key2", 1.23e3, "key3", 'u', "key4", true)
	// Expected output:
	// key=123 key2=1.23e3 key3="u" key4=true

	// You can set permanent pairs as logger context.
	log.With("userID", 1000, "PID", os.GetPID())

	// They will be passed along pairs for each record.
	// log.Log("msg", "details about something")
	// Expect output:
	// userID=1000 PID=12345 msg="details about something"
	
	// You need define even one output: set writer and logging format.
	out:=kiwi.UseOutput(os.StdOut, kiwi.Logfmt)
	
	// Until the output defined log records just saved nowhere.
	// You can define arbitrary number of outputs. Each output has its own set of filters.
	// Filters decide pass or not incoming log record to this output.	
	// Example filters below will pass only records which has key "userID" and has value of level="FATAL".	
	out.With("userID").WithValues("level", "FATAL")
	
	// So in this manner you can fan out log record to several outputs.
	// For example write separate log of critical errors and common log with all errors.
}
```

See more ready to run samples in `cmd` subdirs.


### Thread safety

TBD

### Work with context

TBD

### Evaluating rules of record values

Obsoleted.

TODO Need update!

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
3. No, it is not related to `log15` logger though `kiwi` shares the same logfmt format and some ideas with him.

## Similar works for structured logging

* [logxi](https://github.com/mgutz/logxi)
* [logrus](https://github.com/Sirupsen/logrus)
* [log15](https://github.com/inconshreveable/log15) — another standard No 15 realization :)

## Comparison with other loggers

    $ go test -bench=. -benchmem
    BenchmarkLevelsKiwiTyped-4                100000         18569 ns/op       0.11 MB/s        7991 B/op        115 allocs/op
    BenchmarkLevelsKiwiTypedComplex-4          50000         29985 ns/op       0.07 MB/s       15923 B/op        219 allocs/op
    BenchmarkLevelsKiwi-4                     100000         21562 ns/op       0.09 MB/s        7995 B/op        120 allocs/op
    BenchmarkLevelsKiwiComplex-4               50000         30958 ns/op       0.06 MB/s       16034 B/op        220 allocs/op
    BenchmarkLevelsStdLog-4                   100000         22012 ns/op       0.09 MB/s        7159 B/op        124 allocs/op
    BenchmarkLevelsStdLogComplex-4             50000         32464 ns/op       0.06 MB/s       11446 B/op        200 allocs/op
    BenchmarkLevelsLogxi-4                    100000         13447 ns/op       0.15 MB/s        4127 B/op         74 allocs/op
    BenchmarkLevelsLogxiComplex-4              50000         36282 ns/op       0.06 MB/s       10747 B/op        182 allocs/op
    BenchmarkLevelsLogrus-4                    50000         37246 ns/op       0.05 MB/s       12320 B/op        177 allocs/op
    BenchmarkLevelsLogrusComplex-4             30000         41640 ns/op       0.05 MB/s       13989 B/op        231 allocs/op
    BenchmarkLevelsLog15-4                     30000         54341 ns/op       0.04 MB/s       14998 B/op        224 allocs/op
    BenchmarkLevelsLog15Complex-4              20000         63527 ns/op       0.03 MB/s       18339 B/op        300 allocs/op
	
Well after oftimization kiwi runs faster. It is not the fastest logger among benchmarked but not slowest.
It much faster than `logrus` and `log15` but slower than `logxi`.

## Origins

* logfmt description [brandur.org/logfmt](https://brandur.org/logfmt)
* logfmt realization in Go and specs [godoc.org/github.com/kr/logfmt](https://godoc.org/github.com/kr/logfmt)
* picture used for logo [openclipart.org/detail/4416/kiwi-bird](https://openclipart.org/detail/4416/kiwi-bird)
