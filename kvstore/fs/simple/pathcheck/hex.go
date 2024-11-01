package pchk

type IsHex func(byte) bool

type IsHexL func(byte) bool
type IsHexH func(byte) bool

func (h IsHexH) ToIsHex(l IsHexL) IsHex {
	return func(b byte) bool {
		return h(b) || l(b)
	}
}

func IsHexDefaultL(b byte) bool {
	var h bool = 'a' <= b && b <= 'f'
	return h || IsNumDefault(b)
}
