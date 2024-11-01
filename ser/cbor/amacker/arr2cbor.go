package arr2cbor

import (
	"bytes"
	"context"

	ac "github.com/fxamacker/cbor/v2"

	sc "github.com/takanoriyanagitani/go-cbor2kvstore/ser/cbor"
)

func CborRowToBuffer(
	_ context.Context,
	row sc.CborRow,
	buf *bytes.Buffer,
) error {
	var s []any = row
	return ac.MarshalToBuffer(s, buf)
}

var CborToBufferDefault sc.CborToBuffer = CborRowToBuffer
