package consolelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

type event map[string]interface{}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.TimeFormat == "" {
		w.TimeFormat = defaultTimeFormat
	}

	var buf bytes.Buffer

	var evt event
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	w.writeTimestamp(evt, &buf)
	w.writeLevel(evt, &buf)
	w.writeComponent(evt, &buf)
	w.writeCaller(evt, &buf)

	w.writeMessage(evt, &buf)
	w.writeFields(evt, &buf)

	buf.WriteByte('\n')
	buf.WriteTo(w.Out)
	return len(p), nil
}

func (w ConsoleWriter) writeTimestamp(evt event, buf *bytes.Buffer) {
	var t string
	if tt, ok := evt[zerolog.TimestampFieldName].(string); ok {
		ts, err := time.Parse(time.RFC3339, tt)
		if err != nil {
			t = tt
		} else {
			t = ts.Format(w.TimeFormat)
		}
	}
	fmt.Fprintf(buf, "\x1b[2m%s\x1b[0m", t)
}

func (w ConsoleWriter) writeLevel(evt event, buf *bytes.Buffer) {
	var l string
	var f string
	if ll, ok := evt[zerolog.LevelFieldName].(string); ok {
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
		l = strings.ToUpper(fmt.Sprintf("%s", evt[zerolog.LevelFieldName]))[0:4]
	}
	fmt.Fprintf(buf, " \x1b[%sm%s\x1b[0m", f, l)
}

func (w ConsoleWriter) writeComponent(evt event, buf *bytes.Buffer) {
	var c string
	if cc, ok := evt["component"].(string); ok {
		c = cc
	}
	if len(c) > 0 {
		fmt.Fprintf(buf, " [\x1b[1m%s\x1b[0m]", c)
	}
}

func (w ConsoleWriter) writeCaller(evt event, buf *bytes.Buffer) {
	var c string
	if cc, ok := evt[zerolog.CallerFieldName].(string); ok {
		c = cc
	}
	if len(c) > 0 {
		cwd, err := os.Getwd()
		if err == nil {
			c = strings.TrimPrefix(c, cwd)
			c = strings.TrimPrefix(c, "/")
		}
		fmt.Fprintf(buf, " \x1b[2m\x1b[0m\x1b[1m%s\x1b[2m >\x1b[0m", c)
	}
}

func (w ConsoleWriter) writeMessage(evt event, buf *bytes.Buffer) {
	var m string
	m = fmt.Sprintf("%s", evt[zerolog.MessageFieldName])

	fmt.Fprintf(buf, " \x1b[0m%s\x1b[0m", m)
}

func (w ConsoleWriter) writeFields(evt event, buf *bytes.Buffer) {
	var fields = make([]string, 0, len(evt))
	for field := range evt {
		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName, zerolog.CallerFieldName, "component":
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	if len(fields) > 0 {
		buf.WriteByte(' ')
	}
	for _, field := range fields {
		fmt.Fprintf(buf, "\x1b[2m%s=\x1b[0m", field)
		fmt.Fprintf(buf, "%s ", evt[field])
	}
}
