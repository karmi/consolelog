package main

import (
	"errors"
	"os"

	"github.com/rs/zerolog"

	"github.com/karmi/consolelog/consolelog"
)

func main() {
	output := consolelog.ConsoleWriter{Out: os.Stderr}

	defaultLogger := zerolog.New(output).
		With().
		Timestamp().
		Int("pid", 37556).
		Caller().
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
}
