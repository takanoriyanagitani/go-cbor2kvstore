package pchk

type IsNum func(byte) bool

func IsNumDefault(b byte) bool {
	return '0' <= b && b <= '9'
}
