# NANOLOG

Simple minimalistic logging package with log level capability and colored output

### Basic usage

```go
import (
  log "github.com/stanyx/nanolog"
)

func main() {

  log.Init(log.Options{
    Level: log.DebugLevel
  })

  log.DEBUG.Println("debug message")

}
```

### Known limitations

- colored output not supported for Non-Unix platforms