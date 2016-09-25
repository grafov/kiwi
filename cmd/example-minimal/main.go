package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	kiwi.SinkTo(os.Stdout, kiwi.UseLogfmt()).Start()
	l := kiwi.New()
	l.Add("sample-record", 1).Log()
}
