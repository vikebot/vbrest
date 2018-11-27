package vbcore

import "strconv"

// ItoNumberUnity converts an int to a shortend string representation with
// it's unity. Example: 1000-1999 will be converted to 1K. Works will K, M
// and B. Everything below 1000 will be returned as normal number
func ItoNumberUnity(value int) string {
	if value >= 1000000000 {
		return strconv.Itoa(value/1000000000) + "B"
	} else if value >= 1000000 {
		return strconv.Itoa(value/1000000) + "M"
	} else if value >= 1000 {
		return strconv.Itoa(value/1000) + "K"
	}
	return strconv.Itoa(value)
}
