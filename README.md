A `ConsoleWriter` for <https://github.com/rs/zerolog>.

![Screenshot](screenshot.png)

```go
package main

import (
	"errors"

	"github.com/rs/zerolog"

	"github.com/karmi/consolelog/consolelog"
)

func main() {
	output := consolelog.NewConsoleWriter()

	logger := zerolog.New(output).
		With().
		Timestamp().
		Int("pid", 37556).
		Caller().
		Logger()

	logger.
		Info().
		Str("listen", ":8080").
		Msg("Starting listener")

	logger.
		Error().
		Err(errors.New("connection reset by peer")).
		Str("database", "myapp").
		Str("host", "localhost:4932").
		Msg("Database connection lost")
}
```
