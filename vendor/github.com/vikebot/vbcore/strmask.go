package vbcore

// StrMask replaces the later half of the passed string with *.
func StrMask(str string) string {
	if len(str) == 0 {
		return ""
	} else if len(str) == 1 {
		return "*"
	} else {
		return StrMaskIdx(str, len(str)/2)
	}
}

// StrMaskIdx replaces all characters with * starting at index idx.
func StrMaskIdx(str string, idx int) string {
	res := str[:idx]
	for i := 0; i < len(str)-idx; i++ {
		res += "*"
	}
	return res
}
