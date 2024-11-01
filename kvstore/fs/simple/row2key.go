package fstr

import (
	"context"
	"errors"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
)

type RowToKey kv.RowToKey

type RowToFsKey struct {
	RowToDirName
	RowToFileName
	CreatePath
}

func (f RowToFsKey) ToRowToKey() RowToKey {
	return func(ctx context.Context, row ck.Row) (kv.Key, error) {
		dir, de := f.RowToDirName(ctx, row)
		fil, fe := f.RowToFileName(ctx, row)
		pat, pe := f.CreatePath(ctx, dir, fil)
		return pat, errors.Join(de, fe, pe)
	}
}
