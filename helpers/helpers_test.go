package helpers

import (
	"fmt"
	"slices"
	"testing"
)

func singleTestMyRounded(a int, b int, expected int, t *testing.T) {
	measured := MyRoundedMean(a, b)
	if measured != expected {
		t.Errorf("for a = %d, b = %d -> expected: %d measured: %d\n", a, b, measured, expected)
	}
}

func TestMyRoundedMean(t *testing.T) {
	singleTestMyRounded(0, 4, 2, t)
	singleTestMyRounded(4, 0, 2, t)
	singleTestMyRounded(-4, 0, -2, t)
	singleTestMyRounded(0, -4, -2, t)
	singleTestMyRounded(1, -4, -2, t)
	singleTestMyRounded(-4, 1, -2, t)
	singleTestMyRounded(12, 13, 12, t)
	singleTestMyRounded(-12, -13, -13, t)
}

func singleTestLenToByte(l int, expected byte, t *testing.T) {
	measured := LenToByte(l)
	if measured != expected {
		t.Errorf("With len = %d -> expected '%c' found '%c'", l, expected, measured)
	}
}

func TestLenToByte(t *testing.T) {
	singleTestLenToByte(0, '0', t)
	singleTestLenToByte(1, '1', t)
	singleTestLenToByte(2, '2', t)
	singleTestLenToByte(3, '3', t)
	singleTestLenToByte(4, '4', t)
	singleTestLenToByte(5, '5', t)
	singleTestLenToByte(6, '6', t)
	singleTestLenToByte(7, '7', t)
	singleTestLenToByte(8, '8', t)
	singleTestLenToByte(9, '9', t)
	singleTestLenToByte(10, 'x', t)
	singleTestLenToByte(-1, 'x', t)
	singleTestLenToByte(1234, 'x', t)
	singleTestLenToByte(-1234, 'x', t)
}

func singleTestAll[T interface{}](a []T, f func(t T)bool, expected bool, t *testing.T) {
	measured := All(a, f)
	if measured != expected {
		t.Errorf("for array = %v and f = {...} -> expected: %t measured: %t", a, expected, measured)
	}
}

func TestAll(t *testing.T) {
	singleTestAll([]int{1, 1, 1, 1}, func(i int) bool {return i == 1}, true, t)
	singleTestAll([]int{1, 1, 1, 1}, func(i int) bool {return i != 1}, false, t)
	singleTestAll([]int{1, 2, 4, 8}, func(i int) bool {return i % 2 == 0}, false, t)
	singleTestAll([]int{0, 2, 3, 8}, func(i int) bool {return i % 2 == 0}, false, t)
	singleTestAll([]int{0, 2, 4, 7}, func(i int) bool {return i % 2 == 0}, false, t)
	singleTestAll([]int{0, 2, 4, 8}, func(i int) bool {return i % 2 == 0}, true, t)
}

func singleTestGetOneOnMap[Key comparable, Val comparable](m map[Key]Val, expectedKeys []Key, expectedVals []Val, expectedBool bool, t *testing.T) {
	measuredKey, measuredVal, measuredBool := GetOneOnMap(m)
	if measuredBool != expectedBool || !slices.Contains(expectedKeys, measuredKey) || !slices.Contains(expectedVals, measuredVal) {
		s := fmt.Sprintf("For map = %v -> ", m)
		if measuredBool != expectedBool {
			s += fmt.Sprintf("expectedBool: %t measuredBool: %t ", expectedBool, measuredBool)
		}
		if !slices.Contains(expectedKeys, measuredKey) {
			s += fmt.Sprintf("measuredKey: %v not contained in expectedKeys: %v", measuredKey, expectedKeys)
		}
		if !slices.Contains(expectedVals, measuredVal) {
			s += fmt.Sprintf("measuredVal: %v not contained in expectedVals: %v", measuredVal, expectedVals)
		}
		t.Error(s)
	}
}

func TestGetOneValOnMap(t *testing.T) {
	m1 := map[int]string{
		10: "dd",
	}
	singleTestGetOneOnMap(m1, []int{10}, []string{"dd"}, true, t)
	
	m2 := map[int]string{
		10: "dd",
		20: "ff",
		2: "abc",
	}
	singleTestGetOneOnMap(m2, []int{10, 20, 2}, []string{"ff", "dd", "abc"}, true, t)
	
	m3 := map[float32]byte{}
	singleTestGetOneOnMap(m3, []float32{0}, []byte{0}, false, t)
}