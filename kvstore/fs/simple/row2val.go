package fstr

import (
	"bytes"
	"context"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
)

type RowToVal kv.RowToVal

type RowToValBuf func(context.Context, ck.Row, *bytes.Buffer) error

func (b RowToValBuf) ToRowToVal() RowToVal {
	var buf bytes.Buffer
	return func(ctx context.Context, r ck.Row) (kv.Val, error) {
		buf.Reset()
		e := b(ctx, r, &buf)
		return buf.Bytes(), e
	}
}
