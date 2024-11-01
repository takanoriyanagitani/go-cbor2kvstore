package cbor2kv

import (
	"context"
)

type Row []any

type RowSaver func(context.Context, Row) error
