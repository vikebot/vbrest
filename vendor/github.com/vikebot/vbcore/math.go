package vbcore

// MinInt returns the smaller of the two passed integer
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// MaxInt returns the bigger of the two passed integer
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
