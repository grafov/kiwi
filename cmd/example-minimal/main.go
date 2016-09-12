package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	kiwi.UseOutput(os.Stdout, kiwi.UseLogfmt())
	l := kiwi.New()
	l.Add("sample-record", 1).Log()
}
