package fstr

import (
	"bytes"
	"context"
	"errors"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"

	util "github.com/takanoriyanagitani/go-cbor2kvstore/util"
)

var (
	ErrInvalidIndex error = errors.New("invalid index")
)

type RowToDirName func(context.Context, ck.Row) ([]byte, error)

type RowToDirAny func(context.Context, ck.Row) (any, error)

func SelectColByIdx(_ context.Context, idx uint32, r ck.Row) (any, error) {
	var s []any = r
	var l int = len(s)
	var i int = int(idx)
	if i < l {
		return s[i], nil
	}
	return nil, ErrInvalidIndex
}

func RowToDirAnyFromColIdx(idx uint32) RowToDirAny {
	return util.CurryCtx(SelectColByIdx)(idx)
}

type AnyToDirNameToBuf func(context.Context, any, *bytes.Buffer) error

func (a AnyToDirNameToBuf) ToRowToDirToBuf(r2d RowToDirAny) RowToDirNameToBuf {
	return func(ctx context.Context, r ck.Row, buf *bytes.Buffer) error {
		col, e := r2d(ctx, r)
		if nil != e {
			return e
		}
		return a(ctx, col, buf)
	}
}

type RowToDirNameToBuf func(context.Context, ck.Row, *bytes.Buffer) error

func (b RowToDirNameToBuf) ToRowToDirName() RowToDirName {
	var buf bytes.Buffer
	return func(ctx context.Context, row ck.Row) ([]byte, error) {
		buf.Reset()
		e := b(ctx, row, &buf)
		return buf.Bytes(), e
	}
}
