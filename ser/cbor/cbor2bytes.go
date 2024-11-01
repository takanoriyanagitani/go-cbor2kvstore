package cbor2bytes

import (
	"bytes"
	"context"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
)

type CborRow ck.Row
type Val kv.Val

type CborToBytes func(context.Context, CborRow) (Val, error)

func (b CborToBytes) ToRowToVal() kv.RowToVal {
	return func(ctx context.Context, r ck.Row) (kv.Val, error) {
		v, e := b(ctx, CborRow(r))
		return kv.Val(v), e
	}
}

type CborToBuffer func(context.Context, CborRow, *bytes.Buffer) error

func (b CborToBuffer) ToCborToBytes() CborToBytes {
	var buf bytes.Buffer
	return func(ctx context.Context, row CborRow) (Val, error) {
		buf.Reset()
		e := b(ctx, row, &buf)
		return buf.Bytes(), e
	}
}
