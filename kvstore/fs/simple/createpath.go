package fstr

import (
	"bytes"
	"context"
	"strings"
)

type CreatePath func(ctx context.Context, dir, filename []byte) ([]byte, error)

type BytesToPath func(context.Context, []byte) (string, error)
type BytesToPathBuf func(context.Context, []byte, *strings.Builder) error

func (b BytesToPathBuf) ToBytesToPath() BytesToPath {
	var buf strings.Builder
	return func(ctx context.Context, p []byte) (string, error) {
		buf.Reset()
		e := b(ctx, p, &buf)
		return buf.String(), e
	}
}

func BytesToPathBufDefault(
	_ context.Context,
	b []byte,
	buf *strings.Builder,
) error {
	_, _ = buf.Write(b) // always nil error or OOM
	return nil
}

type PathToBuffer func(
	ctx context.Context,
	dir []byte,
	filename []byte,
	buf *bytes.Buffer,
) error

func (p PathToBuffer) ToCreatePath() CreatePath {
	var buf bytes.Buffer
	return func(ctx context.Context, dir, filename []byte) ([]byte, error) {
		buf.Reset()
		e := p(ctx, dir, filename, &buf)
		return buf.Bytes(), e
	}
}
