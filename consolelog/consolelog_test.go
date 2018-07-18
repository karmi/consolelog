package consolelog_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/karmi/consolelog/consolelog"
)

func TestConsolelog(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		w := consolelog.NewConsoleWriter()

		if w.TimeFormat == "" {
			t.Errorf("Missing w.TimeFormat")
		}

		if w.Formatter("foobar") == nil {
			t.Errorf(`Missing default formatter for "foobar"`)
		}

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		o := w.Formatter("time")(d)
		if o != "12:00AM" {
			t.Errorf(`Unexpected output for date %q: %s`, d, o)
		}
	})

	t.Run("SetFormatter", func(t *testing.T) {
		w := consolelog.NewConsoleWriter()

		w.SetFormatter("time", func(i interface{}) string { return "FOOBAR" })

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		o := w.Formatter("time")(d)
		if o != "FOOBAR" {
			t.Errorf(`Unexpected output from custom "time" formatter: %s`, o)
		}
	})

	t.Run("Write", func(t *testing.T) {
		var out bytes.Buffer
		w := consolelog.NewConsoleWriter()
		w.Out = &out

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		_, err := w.Write([]byte(`{"time" : "` + d + `", "level" : "info", "message" : "Foobar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "12:00AM INF Foobar\n"
		actualOutput := out.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})
}
