# Kiwi logger & context keeper [![Go Report Card](https://goreportcard.com/badge/grafov/kiwi)](https://goreportcard.com/report/grafov/kiwi)

*WIP. API and features is subject of changes. Don't use it!*

![Kiwi bird](flomar-kiwi-bird-300px.png)

*Kiwi* /kiːwiː/ are birds native to New Zealand, in the genus Apteryx and family Apterygidae. They are flightless, have hair-like feathers and smell like mushrooms. They look strange and funny so when I wrote a logger for Go language I decided to devote it to this beast which I never seen in a wild (I live very far from places where kiwis are live).

*Kiwi Logger* — this is a library with an odd logic that log your application data in its own strange way.

## Features offered by structered logging and logfmt generally and by Kiwi particularly

Shortly: both humans and robots will love it!

Briefly:

* structured logging in *logfmt* for high readability by humans
* optional JSON format that liked by machines
* dynamically selectable outputs (changing log verbosity on the fly)
* there are not nailed levels, not hardcoded fields in the format
* can keep context of application

Kiwi logger has two primary design goals:

1. Convenient structured logging syntax: logfmt as default format, method chaining.
2. Separation of logging flow from control flow: you log everything without conditions but output filtering really select what and where will be saved.

Point 2 needs more explanations. 
Traditional way for logging is set a level of severity for each log record. 
And check a level before passing this record to a writer.
It is not bad way but it is not obvious.
Especially when logger introduces many severity levels like these all "debug", "info", "warning", "critical", "fatal", "panic" and so on. 
Look the internet for many guides with controversial recommendations how to distinguish all these "standard" levels and apply them to various events in your application.
When you should use "fatal" instead of "panic" or "debug" instead of "info".

There is alternative simple way not use severity levels at all. Structured logging does things in this way. 
Especially it right for *logfmt* format.
It is not a required part of logfmt or structured logging but it naturally ensued from the rules they offer.
You just log pairs of keys and values and these pairs may be of any kind. There is not a standard list of keys.
If you need you can use levels for log records for example set key named "level" (or any other name you want) with values INFO, WARNING, ERROR etc.
But it is not a requirement.
So you log just pairs of arbitrary keys and values and interprete them as you wish.

Feature of `kiwi` logger is dynamic filtering of incoming records.
You log all data of any severity. These log records passed to all defined outputs (log streams).
And you restrict them by set filters for pass only records and their fields which you want to see in this log stream.
It can be changed in any moment: `kiwi` has methods for filtering by keys, values, ranges of values.
Recipe: export the handler or setup any kind of client for setting these filters in your app.
Then you got ability for dynamically change flow and verbosity of logs. 
For example increase verbosity for a specific module or a single handler and decrease them for the rest of the application.

## Docs [![GoDoc](https://godoc.org/github.com/grafov/kiwi?status.svg)](https://godoc.org/github.com/grafov/kiwi)

See documentation in [the wiki](https://github.com/grafov/kiwi/wiki). 
Examples of logger usage see at [cmd/*](cmd) subfolders.
And of course for API description look at [godoc](http://godoc.org/github.com/grafov/kiwi).

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
	
	// You need define even one output: set writer and logging format.
	out:=kiwi.UseOutput(os.StdOut, kiwi.Logfmt)
	
	// Until the output defined log records just saved nowhere.
	// You can define arbitrary number of outputs. Each output has its own set of filters.
	// Filters decide pass or not incoming log record to this output.	
	// Example filters below will pass only records which has key "userID" and has value of level="FATAL".	
	out.With("userID").WithValues("level", "FATAL")
	
	// So in this manner you can fan out log record to several outputs.
	// For example write separate log of critical errors and common log with all errors.
	// By default without any filters any output accepts any incoming log records.
	out2 := kiwi.UseOutput(os.StdErr, kiwi.JSON)

	// Kiwi offers various filters for set conditions for outputs.
	out2.WithInt64Range("userID", 100, 500).WithoutValues("label", "debug")
}
```

See more ready to run samples in `cmd` subdirs.

## Work with context

`Kiwi` logger has ability keep some pairs during lifetime of a logger instance.

```go
import "github.com/grafov/kiwi"

func main() {
	// Creates a new logger instance.
	log1st := kiwi.New()

	// You can set permanent pairs as logger context.
	log1st.With("userID", 1000, "PID", os.GetPID())

	// They will be passed among other pairs for each record.
	log1st.Log("msg", "details about something")
	// Expect output:
	// userID=1000 PID=12345 msg="details about something"
	
	// Context copied into a new logger instance after logger cloned.
	log2nd := log1st.New()
	
	log2nd.Log("key", "value")
	// Expect output:
	// userID=1000 PID=12345 key="value"
	
	// Get previously keeped context values. Results returned as map[string]interface{}
	appContext := log2nd.GetContext()
	fmt.Printf("%+v\n", appContext)
	
	// You can reset context at any time with
	log2nd.ResetContext()
}
```

## Thread safety

It is unsafe by design. Firstly I have used version for safe work in multiple goroutines.
And it was not only slow but in just not need in many cases. 
If you need a new logger in another execution thread you will create another instanse. Better is clone old instance to a new one for passing the context to a subroutine. It is all.

```go
	// Creates a new logger instance.
	log1st := kiwi.New()

	// Just clone old instance to a new one with keeping of context.
	log2nd := log1st.New()

	// So other concurrent routines may accept logger with the same context.
	go subroutine(log2nd, otherArgs...)
```

For the small apps where you won't init all these instances you would like use global kiwi.Log() method.
This method just immediately flush it's args to an outputs. And by design it is safe for concurrent usage.
Also due design simplicity it not supports context, only regular values. If you need context then you 
application is complex thing hence you will need initialize a new instance of kiwi.Logger().

### Evaluating rules

* Keys and values evaluated *immediately* after they added to a record.
* Context values evaluated *once* when they added to a logger.
* For lazy evaluating of context and record values pass them as functions:

        # For lazy evaluating you need function that returns string
        func longActionForDelayedEvaluation() string {
           // Do something complex...
		   // and got for example integer result.
		   //
		   // You need convert the result to a string.
           return strconv.Itoa(result)
        }
        myLog.Add("lazy-sample", longActionForDelayedEvaluation) # but not longActionForDelayedEvaluation()

Logger accepts functions without args that returns a string: `func () string`.
Hence value of `lazy-sample` from the example above will be evaluated only on `Log()` call.

## Instead of FAQ

0. Kiwi logger not strictly follow logfmt specs.
1. Ideas of key-value format very near to JSON output but with orientation on readability for humans without additional tools for log parsing.
2. Yes, it was architectured and developed to be a standard number 15 that competing with others. It is not pretend to be log format for everything.
3. No, it is not related to `log15` logger though `kiwi` shares the same logfmt format and some ideas with him.
4. It did not offer "Logger" interface because IMO interface is useless for loggers. It offers interfaces for the parts of logger like formatters.

## Similar works for structured logging

* [logxi](https://github.com/mgutz/logxi)
* [logrus](https://github.com/Sirupsen/logrus)
* [log15](https://github.com/inconshreveable/log15) — another standard No 15 realization :)

## Comparison with other loggers

Briefly: it looks not bad :)

    $ go test -bench=. -benchmem
    BenchmarkLevelsKiwiTyped-4                    100000         19538 ns/op       0.10 MB/s        7131 B/op         99 allocs/op
    BenchmarkLevelsKiwiTypedComplex-4              50000         33557 ns/op       0.06 MB/s       12841 B/op        207 allocs/op
    BenchmarkLevelsKiwiTypedHelpers-4             100000         19997 ns/op       0.10 MB/s        6362 B/op         91 allocs/op
    BenchmarkLevelsKiwiTypedHelpersComplex-4       50000         31168 ns/op       0.06 MB/s       12585 B/op        204 allocs/op
    BenchmarkLevelsKiwi-4                         100000         15607 ns/op       0.13 MB/s        7599 B/op        100 allocs/op
    BenchmarkLevelsKiwiComplex-4                   50000         34161 ns/op       0.06 MB/s       13080 B/op        200 allocs/op
    BenchmarkLevelsStdLog-4                       100000         21593 ns/op       0.09 MB/s        7159 B/op        124 allocs/op
    BenchmarkLevelsStdLogComplex-4                 50000         33517 ns/op       0.06 MB/s       11446 B/op        200 allocs/op
    BenchmarkLevelsLogxi-4                        100000         14730 ns/op       0.14 MB/s        4127 B/op         74 allocs/op
    BenchmarkLevelsLogxiComplex-4                  30000         42678 ns/op       0.05 MB/s       10361 B/op        182 allocs/op
    BenchmarkLevelsLogrus-4                        50000         37325 ns/op       0.05 MB/s       12320 B/op        177 allocs/op
    BenchmarkLevelsLogrusComplex-4                 30000         42411 ns/op       0.05 MB/s       13989 B/op        231 allocs/op
    BenchmarkLevelsLog15-4                         30000         53028 ns/op       0.04 MB/s       14993 B/op        224 allocs/op
    BenchmarkLevelsLog15Complex-4                  20000         65809 ns/op       0.03 MB/s       18340 B/op        300 allocs/op
    BenchmarkLevelsGokit-4                        100000         15169 ns/op       0.13 MB/s        2865 B/op         64 allocs/op
    BenchmarkLevelsGokitComplex-4                  30000         34780 ns/op       0.06 MB/s        8486 B/op        164 allocs/op

It is not the fastest logger among benchmarked but fast enough and careful about memory allocations.
It much faster than `logrus` and `log15`. But slower than `logxi` and `gokit` in some tests. Need more detailed tests though.
See the benchmarks in the [github.com/grafov/go-loggers-comparison](https://github.com/grafov/go-loggers-comparison).

## Roadmap

* custom filters
* global thread-safe logger for simple apps
* optional colour formatter for the console
* throttling mode for outputs

## Origins

* logfmt description [brandur.org/logfmt](https://brandur.org/logfmt)
* logfmt realization in Go and specs [godoc.org/github.com/kr/logfmt](https://godoc.org/github.com/kr/logfmt)
* picture used for logo [openclipart.org/detail/4416/kiwi-bird](https://openclipart.org/detail/4416/kiwi-bird)
