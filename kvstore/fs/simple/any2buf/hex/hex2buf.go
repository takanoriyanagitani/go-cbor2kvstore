package hex2buf

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
)

const DirnameDefault string = "untitled.d"
const FilenameDefault string = "untitled.cbor"

type BytesToBuffer func(context.Context, []byte, *bytes.Buffer) error

func BytesToBufferNewStaticFromStr(s string) BytesToBuffer {
	var b []byte = []byte(s)
	return func(_ context.Context, _ []byte, buf *bytes.Buffer) error {
		_, _ = buf.Write(b) // error is always nil or panic
		return nil
	}
}

var BytesToBufferDefaultDirName BytesToBuffer = BytesToBufferNewStaticFromStr(
	DirnameDefault,
)

var BytesToBufferDefaultFile BytesToBuffer = BytesToBufferNewStaticFromStr(
	FilenameDefault,
)

func BytesToBufferHexSimpleNew() BytesToBuffer {
	var hbuf bytes.Buffer
	var enc io.Writer = hex.NewEncoder(&hbuf)
	return func(_ context.Context, b []byte, buf *bytes.Buffer) error {
		hbuf.Reset()

		// writes the b and fill the hbuf with hex
		_, e := enc.Write(b)
		if nil != e {
			return e
		}

		// writes the hex bytes
		_, _ = buf.Write(hbuf.Bytes()) // always nil error or panic
		return nil
	}
}
