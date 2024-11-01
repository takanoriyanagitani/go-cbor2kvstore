package kv

import (
	"context"
	"errors"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
)

type Key []byte
type Val []byte

type RowToKey func(context.Context, ck.Row) (Key, error)

type RowToVal func(context.Context, ck.Row) (Val, error)

type Kvstore func(context.Context, Key, Val) error

func (k Kvstore) ToRowStore(r2k RowToKey, r2v RowToVal) ck.RowSaver {
	return func(ctx context.Context, row ck.Row) error {
		key, ek := r2k(ctx, row)
		val, ev := r2v(ctx, row)
		e := errors.Join(ek, ev)
		if nil != e {
			return e
		}
		return k(ctx, key, val)
	}
}
