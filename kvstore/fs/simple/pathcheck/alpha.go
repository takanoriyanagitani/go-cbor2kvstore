package pchk

type IsAlpha func(byte) bool

type IsAlphaL func(byte) bool
type IsAlphaH func(byte) bool

func (h IsAlphaH) ToIsAlpha(l IsAlphaL) IsAlpha {
	return func(b byte) bool {
		return h(b) || l(b)
	}
}
