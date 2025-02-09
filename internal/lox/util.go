package lox

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_' || c == '-'
}

func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
