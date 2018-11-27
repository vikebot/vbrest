package vbcore

// TernaryOperatorA behaves like a ternary operator on a string
func TernaryOperatorA(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
}

// TernaryOperatorI behaves like a ternary operator on a int
func TernaryOperatorI(condition bool, trueValue int, falseValue int) int {
	if condition {
		return trueValue
	}
	return falseValue
}
