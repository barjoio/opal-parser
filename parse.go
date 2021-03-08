package opalparser2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
)

// Parser is used to parse Opal documents
type Parser struct {
	input      string    // the string input containing markup
	char       charType  // the current character
	frame      string    // the current sliding window selection
	len        int       // the length of the string input
	start      int       // the start of the sliding window
	pos        int       // the end of the sliding window (current position)
	ln         int       // the line number
	col        int       // the column number (position within line)
	startLn    int       // the starting line of a node
	startCol   int       // the starting column of a node
	parseFn    parseFn   // the current parse function
	parseStack []parseFn // a stack of parse functions
	tree       []*Node   // the abstract syntax tree
	nodeStack  []*Node   // a stack of nodes
}

// New is used to create a new parser
func New() *Parser {
	return &Parser{
		pos:      -1,
		ln:       1,
		startCol: 1,
		startLn:  1,
	}
}

// Parse is used to parse a raw string input of Opal markup
func (p *Parser) Parse(input string) {
	p.input = input
	p.len = len(input)
	p.parseFn = parseBegin

	for p.parseFn != nil {
		p.parseFn = p.parseFn(p)
	}

	// add to tree
	p.addToParent()

	fmt.Println("parsing finished.")
	b, err := json.MarshalIndent(p.tree, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

// ParseFile is used to parse files containing Opal markup
func (p *Parser) ParseFile(filein string) {
	b, err := ioutil.ReadFile(filein)
	if err != nil {
		panic(err)
	}
	p.Parse(string(b))
}

// pushParseFn appends the current parse function to the parse function stack
func (p *Parser) pushParseFn() {
	p.parseStack = append(p.parseStack, p.parseFn)
}

// popParseFn returns and removes the topmost parse function from the parse function stack
func (p *Parser) popParseFn() parseFn {
	top := p.parseStack[len(p.parseStack)-1]
	p.parseStack = p.parseStack[:len(p.parseStack)-1]
	return top
}

// flattenFrame brings the start of the frame up the current position
// effectively skipping the frame content as if it were a node
func (p *Parser) flattenFrame() {
	p.start = p.pos
	p.startCol = p.col
	p.startLn = p.ln
	p.frame = p.input[p.start:p.pos]
}

// next is used to advance the parsing state
func (p *Parser) next() {
	p.pos++
	p.frame = p.input[p.start:p.pos]

	if p.pos >= p.len {
		p.char = eof
		return
	}

	p.char = charType(p.input[p.pos])

	p.col++
	if p.char == charNewline {
		p.ln++
		p.col = 1
	}

	if p.isTerminator() {
		p.char = terminator
	}
}

// nextUntil advances the parser until one of the destination options are encountered
func (p *Parser) nextUntil(targetNode nodeType, destinationOptions, invalidChars string) {
	if p.char == eof || p.char == terminator {
		return
	}
	for {
		switch p.char {
		case eof:
			p.addChild(targetNode)
			return
		case terminator:
			p.addChild(targetNode)
			return
		}
		if strings.ContainsRune(destinationOptions, rune(p.char)) {
			p.addChild(targetNode)
			return
		}
		if strings.ContainsRune(invalidChars, rune(p.char)) {
			fmt.Printf("unexpected: %q\n", string(p.char))
			p.addChild(targetNode)
			return
		}
		if p.char == charNewline {
			p.addChild(targetNode)
		}
		p.next()
	}
}

func (p *Parser) skipWhitespace(allowTerminators bool) {
	for {
		if p.char == eof || !unicode.IsSpace(rune(p.char)) && p.char != terminator || !allowTerminators && p.char == terminator {
			return
		}
		p.next()
	}
}

// isTerminator checks for a semicolon or group of whitespace containing 2 or more newlines
func (p *Parser) isTerminator() bool {
	if p.char != charSemicolon && p.char != charNewline {
		return false
	}

	if p.char == charSemicolon {
		return true
	}

	var count int
	for i := p.pos; i < p.len && count < 2 && unicode.IsSpace(rune(p.input[i])); i++ {
		if p.input[i] == '\n' {
			count++
		}
	}
	if count == 2 {
		return true
	}
	return false
}

// trimSpace trims the current frame of whitespace
func (p *Parser) trimSpace() {
	p.frame = strings.TrimSpace(p.frame)
}
