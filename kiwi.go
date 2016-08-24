package kiwi

import "io"

func init() {
	outputs.m = make(map[io.Writer]*Output)
}
