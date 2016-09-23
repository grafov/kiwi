package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	// Bind a new logger to a variable. You may create any number of loggers.
	log := kiwi.New()

	// For starting write log records to some writer output should be initialized.
	output := kiwi.SinkTo(os.Stdout, kiwi.UseLogfmt())

	log.Add("sample-record", 1, "key", "value")
	log.Log()

	// Most logger and output operations support chaining.
	log.Add("sample-record", 2, "key", "value", "key2", 123).Log()

	// On pause output will drop any incoming records.
	output.Pause()
	log.Add("this record will be dropped because single output we declared is on pause")
	output.Continue()

	// You can explicitly remove output but it will automatically closed on application exit.
	output.Close()
}
