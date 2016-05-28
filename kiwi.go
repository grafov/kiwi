package kiwi

import "io"

func init() {
	outputs.w = make(map[io.Writer]*Output)
}
