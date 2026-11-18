package helpers

func MyRoundedMean(a int, b int) int {
	sum := a + b
	if sum < 0 && (-sum) % 2 == 1 {
		return (sum - 1) / 2
	}
	return sum / 2
}

func LenToByte(l int) byte {
	if l >= 0 && l < 10 {
		return '0' + byte(l)
	}
	return 'x'
}

func All[T interface{}](a []T, f func(t T)bool) bool {
	for _, elt := range(a) {
		if !f(elt) {
			return false
		}
	}
	return true
}

func GetOneOnMap[Key comparable, Val interface{}](m map[Key]Val) (Key, Val, bool) {
	for k, v := range m {
		return k, v, true
	}

	var zeroKey Key
	var zeroVal Val
	return zeroKey, zeroVal, false
}
