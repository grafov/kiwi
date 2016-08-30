package main

import (
	"fmt"
	"os"
	"time"

	"github.com/grafov/kiwi"
)

func main() {
	// Bind a new logger to a variable. You may create any number of loggers.
	ctx := kiwi.NewLogger()

	// For starting write ctx records to some writer output should be initialized.
	kiwi.UseOutput(os.Stdout, kiwi.Logfmt)

	// setup context of the logger
	ctx.With("userID", 1000, "host", "local", "startedAt", time.Now())

	// This record will be supplemented by startedAt value of time.Now().String()
	ctx.Add("sample", 1).Log()

	// This record also will be supplemented by the same value of the time.
	// Because context value evalueted when it was added by ctx.With().
	ctx.Add("sample", 2).Log()

	// You can provide deferred evaluation of context or ctx values if you add them wrapped
	// with func() interface{}, where interface should be one of scalar golang types.
	ctx.With("currentTime", func() string { return time.Now().String() })

	// Get previously saved context for use in the application.
	// They were keep as is without conversion to strings.
	currentContext := make(map[interface{}]interface{})
	currentContext = ctx.GetContext()
	fmt.Printf("some of the context values are: %d, %s\n", currentContext["userID"], currentContext["host"])

	// These records will be output each its own currentTime value because currentTime will
	// be evaluated on each Log() call.
	ctx.Add("sample", 3).Log()
	ctx.Add("sample", 4).Log()
	ctx.Flush()
}
