package cbor2kvs

import (
	"context"
	"iter"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"
	ic "github.com/takanoriyanagitani/go-cbor2kvstore/iter/cbor2array"
)

type CborToIterToKvstore struct {
	ic.CborToArrays
	ck.RowSaver
}

func (c CborToIterToKvstore) SaveAll(ctx context.Context) error {
	var i iter.Seq[[]any] = c.CborToArrays()
	for row := range i {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var r ck.Row = row
		e := c.RowSaver(ctx, r)
		if nil != e {
			return e
		}
	}
	return nil
}
