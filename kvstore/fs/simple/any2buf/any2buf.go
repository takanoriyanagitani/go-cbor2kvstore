package any2buf

import (
	"bytes"
	"context"
	"errors"
	"time"

	cf "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple"

	hb "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/any2buf/hex"
	nb "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/any2buf/num"
	sb "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/any2buf/str"
	tb "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/any2buf/time"
)

var (
	ErrInvalidType error = errors.New("invalid type")
)

type AnyToBuffer struct {
	nb.IntegerToBuffer8
	nb.IntegerToBuffer16
	nb.IntegerToBuffer32
	nb.IntegerToBuffer64

	tb.TimeToBuffer

	sb.StringToBuffer
	hb.BytesToBuffer
}

func (a AnyToBuffer) ToBuf(ctx context.Context, i any, b *bytes.Buffer) error {
	switch t := i.(type) {
	case uint8:
		return a.IntegerToBuffer8(ctx, t, b)
	case uint16:
		return a.IntegerToBuffer16(ctx, t, b)
	case uint32:
		return a.IntegerToBuffer32(ctx, t, b)
	case uint64:
		return a.IntegerToBuffer64(ctx, t, b)
	case []byte:
		return a.BytesToBuffer(ctx, t, b)
	case string:
		return a.StringToBuffer(ctx, t, b)
	case time.Time:
		return a.TimeToBuffer(ctx, t, b)
	}
	return ErrInvalidType
}

func (a AnyToBuffer) AsAnyToDirnameToBuf() cf.AnyToDirNameToBuf {
	return a.ToBuf
}

func AnyToBufferDefaultNew() AnyToBuffer {
	return AnyToBuffer{
		IntegerToBuffer8:  nb.IntegerToBuffer8(nb.IntToBufFmt8),
		IntegerToBuffer16: nb.IntegerToBuffer16(nb.IntToBufFmt16),
		IntegerToBuffer32: nb.IntegerToBuffer32(nb.IntToBufFmt32),
		IntegerToBuffer64: nb.IntegerToBuffer64(nb.IntToBufFmt64),

		TimeToBuffer: tb.TimeToBufYyyyMmDd,

		StringToBuffer: sb.StringToBuf,
		BytesToBuffer:  hb.BytesToBufferHexSimpleNew(),
	}
}
