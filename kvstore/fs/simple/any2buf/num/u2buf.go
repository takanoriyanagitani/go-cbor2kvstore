package int2buf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidInteger error = errors.New("invalid integer")
)

type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type IntegerToBuffer[T Unsigned] func(context.Context, T, *bytes.Buffer) error

type IntegerToBuffer8 IntegerToBuffer[uint8]
type IntegerToBuffer16 IntegerToBuffer[uint16]
type IntegerToBuffer32 IntegerToBuffer[uint32]
type IntegerToBuffer64 IntegerToBuffer[uint64]

func IntToBufFmtNew[T Unsigned](f string) IntegerToBuffer[T] {
	return func(_ context.Context, i T, buf *bytes.Buffer) error {
		_, e := fmt.Fprintf(buf, f, i)
		return e
	}
}

var IntToBufFmt8 IntegerToBuffer[uint8] = IntToBufFmtNew[uint8]("%02x")
var IntToBufFmt16 IntegerToBuffer[uint16] = IntToBufFmtNew[uint16]("%04x")
var IntToBufFmt32 IntegerToBuffer[uint32] = IntToBufFmtNew[uint32]("%08x")
var IntToBufFmt64 IntegerToBuffer[uint64] = IntToBufFmtNew[uint64]("%016x")
