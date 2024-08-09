package gls

//go:noinline
func Switcher[A any, B any](msg string, code uint64, fn func(string) (A, B)) (A, B) {
	c := code % 16
	r := code / 16
	if r == 0 {
		return fn(msg)
	}
	switch c {
	case 0:
		return Switcher(msg+"n0", r, fn)
	case 1:
		return Switcher(msg+"n1", r, fn)
	case 2:
		return Switcher(msg+"n2", r, fn)
	case 3:
		return Switcher(msg+"n3", r, fn)
	case 4:
		return Switcher(msg+"n4", r, fn)
	case 5:
		return Switcher(msg+"n5", r, fn)
	case 6:
		return Switcher(msg+"n6", r, fn)
	case 7:
		return Switcher(msg+"n7", r, fn)
	case 8:
		return Switcher(msg+"n8", r, fn)
	case 9:
		return Switcher(msg+"n9", r, fn)
	case 10:
		return Switcher(msg+"n10", r, fn)
	case 11:
		return Switcher(msg+"n11", r, fn)
	case 12:
		return Switcher(msg+"n12", r, fn)
	case 13:
		return Switcher(msg+"n13", r, fn)
	case 14:
		return Switcher(msg+"n14", r, fn)
	case 15:
		return Switcher(msg+"n15", r, fn)
	}

	return fn("")
}
