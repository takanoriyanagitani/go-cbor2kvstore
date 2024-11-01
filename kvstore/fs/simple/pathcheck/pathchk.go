package pchk

import (
	"bytes"
	"context"
	"errors"

	kf "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple"
)

type PathSeparator byte

const PathSeparatorDefault PathSeparator = '/'

type DirnameChecker func(context.Context, []byte) error
type FilenameChecker func(context.Context, []byte) error

type PathChecker struct {
	DirnameChecker
	FilenameChecker
	PathSeparator
}

func (p PathChecker) ToPathToBuffer() kf.PathToBuffer {
	return func(
		ctx context.Context,
		dir []byte,
		filename []byte,
		buf *bytes.Buffer,
	) error {
		var de error = p.DirnameChecker(ctx, dir)
		var fe error = p.FilenameChecker(ctx, filename)
		var e error = errors.Join(de, fe)
		if nil != e {
			return e
		}

		_, _ = buf.Write(dir) // always nil error or panic
		_, _ = buf.Write([]byte{byte(p.PathSeparator)})
		_, _ = buf.Write(filename) // always nil error or panic
		return nil
	}
}

var PathCheckerDefault PathChecker = PathChecker{
	DirnameChecker:  ByteChecker(IsHexDefaultL).AsDirnameChecker(),
	FilenameChecker: ByteChecker(IsHexDefaultL).AsFilenameChecker(),
	PathSeparator:   PathSeparatorDefault,
}
