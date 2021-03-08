package opalparser2

type charType rune

// list of characters found within the grammar of the markup
const (
	eof           charType = -1
	terminator    charType = -2
	charFullstop  charType = '.'
	charColon     charType = ':'
	charSlash     charType = '/'
	charNewline   charType = '\n'
	charSemicolon charType = ';'
	charGrave     charType = '`'
	charHyphen    charType = '-'
)

var charsWhitespace = "\t\n\v\f\r " + string(rune(0x85)) + string(rune(0xA0))
