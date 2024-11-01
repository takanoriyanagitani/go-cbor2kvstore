package time2buf

import (
	"bytes"
	"context"
	"time"
)

type TimeToBuffer func(context.Context, time.Time, *bytes.Buffer) error

func TimeToBufFromFmt(layout string) TimeToBuffer {
	return func(_ context.Context, t time.Time, buf *bytes.Buffer) error {
		var s string = t.Format(layout)
		_, _ = buf.WriteString(s) // error is always nil or panic
		return nil
	}
}

var TimeToBufRfc3339 TimeToBuffer = TimeToBufFromFmt(time.RFC3339)
var TimeToBufDate TimeToBuffer = TimeToBufFromFmt(time.DateOnly)
var TimeToBufYyyyMmDd TimeToBuffer = TimeToBufFromFmt("20060102")
