# Kiwi logger & context keeper [![Go Report Card](https://goreportcard.com/badge/grafov/kiwi)](https://goreportcard.com/report/grafov/kiwi)

*unfinished, development in progress*

This software pretends to be an odd piece of code.

![Kiwi bird](flomar-kiwi-bird-300px.png)

*Kiwi* /kiːwiː/ are birds native to New Zealand, in the genus Apteryx and family Apterygidae. They are flightless, have hair-like feathers and smell like a mushrooms. They look strange and funny so when I wrote a logger for Go language I decided to devote it to this beast which I never seen in a wild (I live very far from places where kiwis live).

## Features

* priority on high readability for humans
* JSON mode that liked by machines
* different outputs for log record selected with filters
* no nailed levels, no hardcoded fields, but explicit labels for each logged value
* dynamically change fields visibility in an output
* keeps a context of an application

## Docs [![GoDoc](https://godoc.org/github.com/grafov/kiwi?status.svg)](https://godoc.org/github.com/grafov/kiwi)

Examples of logger usage see at [cmd/*](cmd) subfolders.
See API description and code samples in [godoc](http://godoc.org/github.com/grafov/kiwi).


## Instead of FAQ

0. Kiwi logger not strictly follow logfmt specs.
1. Ideas of key-value format very near to JSON output but with orientation on readability for humans without additional tools for log parsing.
2. Yes, it was architectured and developed to be a standard number 15 that competing with others. It is not pretend to be log format for everything.

## Origins

* logfmt description [brandur.org/logfmt](https://brandur.org/logfmt)
* logfmt realization in Go and specs [godoc.org/github.com/kr/logfmt](https://godoc.org/github.com/kr/logfmt)
* picture used for logo [openclipart.org/detail/4416/kiwi-bird](https://openclipart.org/detail/4416/kiwi-bird)
