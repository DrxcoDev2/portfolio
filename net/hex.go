package http_util

func DecConvertHEX(n int) string {
	if n == 0 {
		return "0"
	}

	hexChars := "0123456789abcdef"
	result := ""

	for n > 0 {
		remainder := n % 16
		result = string(hexChars[remainder]) + result
		n /= 16
	}

	return result
}
