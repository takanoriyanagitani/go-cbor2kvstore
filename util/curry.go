package util

import (
	"context"
)

func CurryCtx[T, U, V any](
	f func(context.Context, T, U) (V, error),
) func(T) func(context.Context, U) (V, error) {
	return func(t T) func(context.Context, U) (V, error) {
		return func(ctx context.Context, u U) (V, error) {
			return f(ctx, t, u)
		}
	}
}

func Curry[T, U, V any](
	f func(T, U) V,
) func(T) func(U) V {
	return func(t T) func(U) V {
		return func(u U) V {
			return f(t, u)
		}
	}
}
