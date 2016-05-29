package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	// Bind a new logger to a variable. You may create any number of loggers.
	log := kiwi.NewLogger()

	tmpFile, _ := os.Create("/tmp/something-completely-different.log")

	// You can get arbitrary number of outputs.
	info := kiwi.GetOutput(os.Stdout, kiwi.Logfmt)
	errors := kiwi.GetOutput(os.Stderr, kiwi.Logfmt)
	something := kiwi.GetOutput(tmpFile, kiwi.Logfmt)

	// Each record by default will copied to all outputs.
	log.Add("level", "INFO", "sample-record", 1, "key", "value")
	log.Add("level", "INFO", "sample-record", 2, "something").Log()
	log.Add("level", "ERROR", "msg", "Error description.").Log()
	log.Add("level", "FATAL").Log()

	// Each output allows fitler out some records and write some other.
	// You specify fitler for keys (key filter).
	// Each of these keys should be presented in record.
	errors.With("error", "msg")
	// Also filter may take into account key values. So only records with levels
	// ERROR and FATAL will be passed filter and written to stderr.
	errors.WithSet("level", "ERROR", "FATAL")

	// Vice versa you can filter out some keys.
	info.Without("error")
	// And define another set of key-val pairs for distinguish outputs.
	info.WithSet("level", "INFO", "WARNING")

	// It will output all records from outputs above if they have key "something".
	// So you can duplicate some records to several log files based on some criteria.
	something.With("something")

	// So if you not define any clauses (With/Without/WithValues/WithoutValues)
	// then all records will copied to an output.

	// Until you call Log() records not copied to outputs.
	log.Log()
}
