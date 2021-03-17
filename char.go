package opalparser

// list of characters found within the grammar of the markup
const (
	eof           rune = -1
	terminator    rune = -2
	charFullstop  rune = '.'
	charColon     rune = ':'
	charSlash     rune = '/'
	charNewline   rune = '\n'
	charSemicolon rune = ';'
	charGrave     rune = '`'
	charHyphen    rune = '-'
	charBackslash rune = '\\'
)

var charsWhitespace = "\t\n\v\f\r " + string(rune(0x85)) + string(rune(0xA0))
