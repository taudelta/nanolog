# NANOLOG

Simple minimalistic logging package with log level capability and colored output

### Basic usage

```go
import (
  "os"
  log "github.com/stanyx/nanolog"
)

func main() {

  log.Init(log.Options{
    Level: log.DebugLevel,
    // example of overriding default writer
    Debug: log.LoggerOptions{
      Writer: os.Stdout,
    }
  })

  log.Debug().Println("debug message")

}
```

### Known limitations

- colored output not supported for Non-Unix platforms