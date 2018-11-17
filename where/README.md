# Runtime package for Kiwi log

The package adds runtime information (package name, function, line
number) as the keys and the values to the logger context.

```go
import (
  "os"
  
  "github.com/grafov/kiwi"
  "github.com/grafov/kiwi/where"
)

func main() {
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()

	kiwi.With(where.What(where.File | where.Line | where.Func))
	kiwi.Log("key", "value")
}
```

The result log record will be like that:

     lineno=11 file="path/to/main.go" function="main.main" key="value"
