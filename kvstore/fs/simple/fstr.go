package fstr

import (
	"context"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
)

type Path string

type RawStore func(context.Context, Path, kv.Val) error

func (r RawStore) ToKvStore(b2p BytesToPath) kv.Kvstore {
	return func(ctx context.Context, k kv.Key, v kv.Val) error {
		p, e := b2p(ctx, k)
		if nil != e {
			return e
		}
		return r(ctx, Path(p), v)
	}
}

type FsStore struct {
	kv.Kvstore
	RowToKey
	RowToVal
}

func (f FsStore) ToRowSaver() ck.RowSaver {
	return f.Kvstore.ToRowStore(
		kv.RowToKey(f.RowToKey),
		kv.RowToVal(f.RowToVal),
	)
}
