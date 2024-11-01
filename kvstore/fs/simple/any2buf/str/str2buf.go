package str2buf

import (
	"bytes"
	"context"
	"errors"
)

var (
	ErrNotString error = errors.New("invalid string")
)

type StringToBuffer func(context.Context, string, *bytes.Buffer) error

func StringToBuf(_ context.Context, s string, buf *bytes.Buffer) error {
	_, _ = buf.WriteString(s) // error is always nil or panic
	return nil
}
