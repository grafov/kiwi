package main

import (
	"os"

	"github.com/grafov/kiwi"
)

func main() {
	// Bind a new logger to a variable. You may create any number of loggers.
	log := kiwi.New()

	tmpFile, _ := os.Create("/tmp/something-completely-different.log")

	// You can set arbitrary number of outputs.
	// But they will remain unused until you explicitly start them with Start().
	info := kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt())
	errors := kiwi.SinkTo(os.Stderr, kiwi.AsLogfmt())
	something := kiwi.SinkTo(tmpFile, kiwi.AsLogfmt())

	// Each record by default will copied to all outputs.
	// But until you Start() any output the records will just dropped as the sample record below.
	log.Add("just something that will lost")

	// Each output allows filter out any records and write any other.
	// You specify filter for the keys (key filter).
	// Each of these keys should be presented in the record.
	errors.WithKey("error", "msg")
	// The filter may take into account key values. So only records with levels
	// ERROR and FATAL will be passed filter and written to stderr.
	errors.WithValue("level", "ERROR", "FATAL").Start()

	// Vice versa you can filter out some keys.
	info.WithoutKey("error")
	// And define another set of key-val pairs for distinguish outputs.
	info.WithValue("level", "INFO", "WARNING").Start()

	// It will output all records from outputs above if they have key "something".
	// So you can duplicate some records to several log files based on some criteria.
	something.WithKey("something").Start()

	// So if you not define any clauses (WithKey/WithoutKey/WithValue/WithoutValues)
	// then all records will copied to an output.

	// Let's go!
	log.Add("level", "INFO", "sample-record", 1, "key", "value")
	log.Add("level", "INFO", "sample-record", 2, "something").Log()
	log.Add("level", "ERROR", "msg", "Error description.").Log()
	log.Add("level", "FATAL").Log()

	// Until you call Log() records not copied to outputs.
	log.Log()
}
