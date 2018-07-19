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

### Custom configuration

```go
customOutput := consolelog.NewConsoleWriter(
		func(w *consolelog.ConsoleWriter) {
			w.PartsOrder = []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			}
		},
		func(w *consolelog.ConsoleWriter) {
			w.TimeFormat = time.RFC822
		},
		func(w *consolelog.ConsoleWriter) {
			w.SetFormatter(
				zerolog.CallerFieldName,
				func(i interface{}) string { return fmt.Sprintf("%s", i) })
			w.SetFormatter(
				zerolog.LevelFieldName,
				func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%-5s", i)) })
		},
	)

	customLogger := zerolog.New(customOutput).
		With().
		Timestamp().
		Caller().
		Logger()
	customLogger.
		Info().
		Str("foo", "bar").
		Msg("Custom message")
```
