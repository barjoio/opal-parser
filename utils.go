package opalparser

import "unicode"

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}
