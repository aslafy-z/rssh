package utils

func Min(x, y uint16) uint16 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}
