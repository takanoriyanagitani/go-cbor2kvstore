package util

func Compose[T, U, V any](
	f func(T) U,
	g func(U) V,
) func(T) V {
	return func(t T) V {
		var u U = f(t)
		return g(u)
	}
}
