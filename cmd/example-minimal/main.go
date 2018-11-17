package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	kiwi.Log("key1", "text value", "key2", 123)
}
