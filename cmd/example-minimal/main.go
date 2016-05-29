package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	kiwi.GetOutput(os.Stdout, kiwi.Logfmt)
	l := kiwi.NewLogger()
	l.Add("sample-record", 1).Log()
}
