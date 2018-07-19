package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/karmi/consolelog/consolelog"
)

func main() {
	output := consolelog.NewConsoleWriter()

	defaultLogger := zerolog.New(output).
		With().
		Timestamp().
		Int("pid", 37556).
		// Caller().
		Logger()

	// https://raw.githubusercontent.com/rs/zerolog/master/pretty.png

	defaultLogger.
		Info().
		Str("listen", ":8080").
		Msg("Starting listener")

	defaultLogger.
		Debug().
		Str("database", "myapp").
		Str("host", "localhost:4932").
		Msg("Connecting to DB")

	defaultLogger.
		Info().
		Str("method", "GET").
		Str("path", "/users").
		Int("resp_time", 23).
		Msg("Access")

	defaultLogger.
		Info().
		Str("method", "POST").
		Str("path", "/posts").
		Int("resp_time", 532).
		Msg("Access")

	defaultLogger.
		Warn().
		Str("method", "POST").
		Str("path", "/posts").
		Int("resp_time", 532).
		Msg("Slow request")

	defaultLogger.
		Info().
		Str("method", "GET").
		Str("path", "/users").
		Int("resp_time", 10).
		Msg("Access")

	defaultLogger.
		Error().
		Err(errors.New("connection reset by peer")).
		Str("database", "myapp").
		Str("host", "localhost:4932").
		Msg("Database connection lost")

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
}
