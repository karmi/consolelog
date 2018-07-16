package consolelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultTimeFormat = time.Kitchen
)

// ConsoleWriter parses the JSON input and writes an ANSI-colorized output to Out.
//
// It is adapted from the original zerolog.ConsoleWriter;
// see: https://github.com/rs/zerolog/blob/master/console.go
//
type ConsoleWriter struct {
	Out        io.Writer
	TimeFormat string
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.TimeFormat == "" {
		w.TimeFormat = defaultTimeFormat
	}

	var buf bytes.Buffer

	var event map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&event)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	var t string
	if tt, ok := event[zerolog.TimestampFieldName].(string); ok {
		ts, err := time.Parse(time.RFC3339, tt)
		if err != nil {
			// fmt.Printf("ERROR: %s\n", err)
			t = tt
		} else {
			t = ts.Format(w.TimeFormat)
		}
	}
	fmt.Fprintf(&buf, "\x1b[2m%s\x1b[0m", t)

	var l string
	var f string
	if ll, ok := event[zerolog.LevelFieldName].(string); ok {
		switch ll {
		case "debug":
			f = "33"
			l = "DBG"
		case "info":
			f = "32"
			l = "INF"
		case "warn":
			f = "31"
			l = "WRN"
		case "error", "fatal", "panic":
			f = "31;1"
			l = "ERR"
		default:
			f = "0"
		}
	} else {
		l = strings.ToUpper(fmt.Sprintf("%s", event[zerolog.LevelFieldName]))[0:4]
	}
	fmt.Fprintf(&buf, " \x1b[%sm%s\x1b[0m", f, l)

	var c string
	if cc, ok := event["component"].(string); ok {
		c = cc
	}
	if len(c) > 0 {
		fmt.Fprintf(&buf, " [\x1b[1m%s\x1b[0m]", c)
	}

	var m string
	m = fmt.Sprintf("%s", event[zerolog.MessageFieldName])

	fmt.Fprintf(&buf, " \x1b[0m%s\x1b[0m ", m)

	var fields = make([]string, 0, len(event))
	for field := range event {
		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName, "component", "build":
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		fmt.Fprintf(&buf, "\x1b[2m%s=\x1b[0m", field)
		fmt.Fprintf(&buf, "%s ", event[field])
	}

	buf.WriteByte('\n')
	buf.WriteTo(w.Out)
	return len(p), nil
}
