package opalparser2

import "fmt"

type errType string

const (
	errUnexpectedTerm = "Unexpected terminator"
	errUnexpectedEOF  = "Unexpected end of file"
	errUnexpectedChar = "Unexpected character"
	errInvalidTagName = "Invalid tag name"
	errNoURL          = "No URL provided in link tag"
)

func (p *Parser) addError(e errType) {
	var err string
	if p.filepath == "" {
		err = fmt.Sprintf("%s at line %d, column %d", e, p.startLn, p.startCol)
	} else {
		err = fmt.Sprintf("%s at %s:%d:%d", e, p.filepath, p.startLn, p.startCol)
	}
	p.nodeStack[0].Errors = append(p.nodeStack[0].Errors, errType(err))
}

func (p *Parser) addErrorUnexpected() {
	switch p.char {
	case eof:
		p.flattenFrame()
		p.addError(errUnexpectedEOF)
	case terminator:
		p.addError(errUnexpectedTerm)
		p.nextFlat()
	default:
		p.flattenFrame()
		p.addError(errType(fmt.Sprintf("%v %s", errUnexpectedChar, string(p.char))))
		p.nextFlat()
	}
	p.popNode()
}
