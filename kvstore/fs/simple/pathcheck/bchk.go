package pchk

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrPathCheck error = errors.New("path check error")
)

type ByteChecker func(byte) (ok bool)

func (c ByteChecker) CheckBytes(_ context.Context, b []byte) error {
	for _, item := range b {
		var ok bool = c(item)
		var ng bool = !ok
		if ng {
			return fmt.Errorf("%w: rejected byte=%v", ErrPathCheck, item)
		}
	}
	return nil
}

func (c ByteChecker) AsDirnameChecker() DirnameChecker   { return c.CheckBytes }
func (c ByteChecker) AsFilenameChecker() FilenameChecker { return c.CheckBytes }
