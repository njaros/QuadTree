package helpers

func MyRoundedMean(a int, b int) int {
	sum := a + b
	if sum < 0 && (-sum) % 2 == 1 {
		return (sum - 1) / 2
	}
	return sum / 2
}

func LenToByte(l int) byte {
	if l > 0 && l < 10 {
		return '0' + byte(l)
	}
	return 'x'
}