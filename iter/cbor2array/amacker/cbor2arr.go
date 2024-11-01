package cbor2iter

import (
	"io"
	"iter"

	ac "github.com/fxamacker/cbor/v2"

	ic "github.com/takanoriyanagitani/go-cbor2kvstore/iter/cbor2array"
)

type CborToArrayIter struct {
	*ac.Decoder
}

func (c CborToArrayIter) ToIter() iter.Seq[[]any] {
	return func(yield func([]any) bool) {
		var buf []any
		var err error
		for {
			clear(buf)
			buf = buf[:0]

			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArrayIter) AsCborToArrays() ic.CborToArrays {
	return c.ToIter
}

func CborToArrayIterNew(rdr io.Reader) CborToArrayIter {
	return CborToArrayIter{Decoder: ac.NewDecoder(rdr)}
}
