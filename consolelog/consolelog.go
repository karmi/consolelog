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

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

const (
	defaultTimeFormat = time.Kitchen
)

var (
	bold   = color.New(color.Bold).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	faint  = color.New(color.Faint).SprintFunc()

	defaultOutput    = os.Stderr
	defaultFormatter = func(i interface{}) string { return fmt.Sprintf("%s", i) }
)

// ConsoleWriter parses the JSON input and writes an ANSI-colorized output to Out.
//
// It is adapted from the zerolog.ConsoleWriter;
// see: https://github.com/rs/zerolog/blob/master/console.go.
//
type ConsoleWriter struct {
	Out        io.Writer
	TimeFormat string
	formatters map[string]Formatter
}

// Formatter transforms the input into a string.
//
type Formatter func(interface{}) string

type event map[string]interface{}

// NewConsoleWriter creates and initializes a new ConsoleWriter.
//
func NewConsoleWriter() ConsoleWriter {
	w := ConsoleWriter{Out: defaultOutput, TimeFormat: defaultTimeFormat}
	w.formatters = make(map[string]Formatter)

	w.setDefaultFormatters()

	return w
}

// Formatter returns a formatter by id or the default formatter if none is found.
//
func (w ConsoleWriter) Formatter(id string) Formatter {
	if f, ok := w.formatters[id]; ok {
		return f
	}
	return defaultFormatter
}

// SetFormatter registers a formatter function by id.
//
func (w ConsoleWriter) SetFormatter(id string, f Formatter) {
	w.formatters[id] = f
}

// Write appends the output to Out.
//
func (w ConsoleWriter) Write(p []byte) (n int, err error) {
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
	buf.WriteString(w.Formatter("timestamp")(evt[zerolog.TimestampFieldName]))
}

func (w ConsoleWriter) writeLevel(evt event, buf *bytes.Buffer) {
	var l string
	if ll, ok := evt[zerolog.LevelFieldName].(string); ok {
		switch ll {
		case "debug":
			l = yellow("DBG")
		case "info":
			l = green("INF")
		case "warn":
			l = red("WRN")
		case "error":
			l = bold(red("ERR"))
		case "fatal":
			l = bold(red("FTL"))
		case "panic":
			l = bold(red("PNC"))
		default:
			l = bold("N/A")
		}
	} else {
		l = strings.ToUpper(fmt.Sprintf("%s", evt[zerolog.LevelFieldName]))[0:3]
	}
	buf.WriteByte(' ')
	buf.WriteString(l)
}

func (w ConsoleWriter) writeComponent(evt event, buf *bytes.Buffer) {
	var c string
	if cc, ok := evt["component"].(string); ok {
		c = cc
	}
	if len(c) > 0 {
		buf.WriteByte(' ')
		buf.WriteString("[")
		buf.WriteString(bold(c))
		buf.WriteString("]")
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
		buf.WriteByte(' ')
		buf.WriteString(bold(c))
		buf.WriteString(faint(" >"))
	}
}

func (w ConsoleWriter) writeMessage(evt event, buf *bytes.Buffer) {
	var m string
	m = fmt.Sprintf("%s", evt[zerolog.MessageFieldName])
	buf.WriteByte(' ')
	buf.WriteString(m)
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
		buf.WriteString(faint(field))
		buf.WriteString(faint("="))
		fmt.Fprintf(buf, "%s ", evt[field])
	}
}

func (w ConsoleWriter) setDefaultFormatters() {
	w.SetFormatter(
		"timestamp", func(i interface{}) string {
			var t string
			if tt, ok := i.(string); ok {
				ts, err := time.Parse(time.RFC3339, tt)
				if err != nil {
					t = tt
				} else {
					t = ts.Format(w.TimeFormat)
				}
			}
			return faint(t)
		})
}
