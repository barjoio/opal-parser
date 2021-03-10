package opalparser2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Parser is used to parse Opal documents
type Parser struct {
	input      string    // the string input containing markup
	filepath   string    // the file path to the markup file
	char       charType  // the current character
	charWidth  int       // the number of bytes used by the current character
	frame      string    // the current sliding window selection
	len        int       // the length of the string input
	start      int       // the start of the sliding window
	pos        int       // the end of the sliding window (current position)
	ln         int       // the line number
	col        int       // the column number (position within line)
	startLn    int       // the starting line of a node
	startCol   int       // the starting column of a node
	firstSpace charType  // stores the first encountered space in a set of whitespace
	linesHere  int       // stores the number of lines encountered through a set of whitespace
	parseFn    parseFn   // the current parse function
	parseStack []parseFn // a stack of parse functions
	tree       []*Node   // the abstract syntax tree
	nodeStack  []*Node   // a stack of nodes
	nodeType   nodeType  // the nodeType of the current parent node
}

// New is used to create a new parser
func New() *Parser {
	return &Parser{
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

	p.createNode(nodeRoot)

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
	p.filepath = filein
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
	p.frame = ""
}

// next is used to advance the parsing state
func (p *Parser) next() {
	// add current character to frame
	if p.char != 0 {
		p.frame += string(p.char)
	}

repeat:

	// check for eof
	if p.pos >= p.len-1 {
		p.char = eof
		p.pos++
		return
	}

	// increment pos/col
	p.pos += p.charWidth
	p.col++

	// get new char
	r, w := utf8.DecodeRuneInString(p.input[p.pos:])
	p.char = charType(r)
	p.charWidth = w

	// handle newlines
	if p.char == charNewline {
		p.ln++
		p.col = 0
		p.linesHere++
	}

	// skip repeating whitespace
	if unicode.IsSpace(rune(p.char)) {
		// store first encountered whitespace to set as char afterwards
		if p.firstSpace == 0 {
			p.firstSpace = p.char
		}
		// don't add spaces to an empty frame
		if len(p.frame) == 0 {
			goto repeat
		}
		// lookahead by 1 to keep skipping space
		if p.pos+1 < p.len && unicode.IsSpace(rune(p.input[p.pos+1])) {
			goto repeat
		}
		// if 2 or more newlines have been encountered, set char to terminator
		if p.linesHere >= 2 {
			p.char = terminator
		} else {
			p.char = p.firstSpace
		}
	}
	p.firstSpace, p.linesHere = 0, 0
}

// nextFlat calls `next` then flattens the frame
func (p *Parser) nextFlat() {
	p.next()
	p.flattenFrame()
}

// nextUntil advances the parser until one of the destination options are encountered
func (p *Parser) nextUntil(destinationOptions string) {
	if p.char == eof || p.char == terminator {
		return
	}
	for {
		switch p.char {
		case eof, terminator:
			return
		}
		if strings.ContainsRune(destinationOptions, rune(p.char)) {
			return
		}
		p.next()
	}
}

func (p *Parser) skipWhitespace() {
	if unicode.IsSpace(rune(p.char)) {
		p.next()
	}
}

func (p *Parser) skipTo(destination charType) {
	for {
		switch p.char {
		case eof, terminator:
			return
		case destination:
			p.nextFlat()
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

func (p *Parser) trimFrame() {
	p.frame = strings.TrimSpace(p.frame)
}

func (p *Parser) debug() {
	for {
		fmt.Printf("%q,%d,%q,%d\n", p.frame, p.pos, string(p.char), p.char)
		if p.char == eof {
			os.Exit(1)
		}
		p.next()
	}
}
