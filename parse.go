package opalparser

import (
	"io/ioutil"
	"strings"
	"unicode"
)

// Parser is used to parse Opal documents
type Parser struct {
	input      []rune   // the string input containing markup
	filepath   string   // the file path to the markup file
	char       rune     // the current character
	frame      []rune   // the current sliding window selection
	len        int      // the length of the string input
	start      int      // the start of the sliding window
	pos        int      // the end of the sliding window (current position)
	ln         int      // the line number
	col        int      // the column number (position within line)
	startLn    int      // the starting line of a node
	startCol   int      // the starting column of a node
	firstSpace rune     // stores the first encountered space in a set of whitespace
	linesHere  int      // stores the number of lines encountered through a set of whitespace
	ignoreChar bool     // switch to deciding if the current character should be ignored
	parseFn    parseFn  // the current parse function
	tree       []*Node  // the abstract syntax tree
	nodeStack  []*Node  // a stack of nodes
	nodeType   nodeType // the nodeType of the current parent node
}

// New is used to create a new parser
func New() *Parser {
	return &Parser{
		pos:      -1,
		ln:       1,
		col:      0,
		startCol: 1,
		startLn:  1,
	}
}

// Parse is used to parse a raw string input of Opal markup
func (p *Parser) Parse(input string) {
	p.input = []rune(input)
	p.len = len(p.input)
	p.parseFn = parseBegin

	p.createNode(nodeRoot)
	p.next()

	for p.parseFn != nil {
		p.parseFn = p.parseFn(p)
	}

	// add to tree
	p.addToParent()
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

// flattenFrame brings the start of the frame up the current position
// effectively skipping the frame content as if it were a node
func (p *Parser) flattenFrame() {
	p.start = p.pos
	p.startCol = p.col
	p.startLn = p.ln
	p.frame = []rune{}
}

// next is used to advance the parsing state
func (p *Parser) next() {
	// add current character to frame
	if p.char != 0 && !p.ignoreChar {
		p.frame = append(p.frame, p.char)
	}

repeat:

	// check for eof
	if p.pos >= p.len-1 {
		p.char = eof
		p.pos++
		return
	}

	// increment pos/col
	p.pos++
	p.col++

	// get new char
	p.char = p.input[p.pos]

	// handle newlines
	if p.char == charNewline {
		p.ln++
		p.col = 0
		p.linesHere++
	}

	// check for terminator
	if p.char == charSemicolon {
		p.char = terminator
		return
	}

	// skip repeating whitespace
	if unicode.IsSpace(p.char) {
		// store first encountered whitespace to set as char afterwards
		if p.firstSpace == 0 {
			p.firstSpace = p.char
		}
		if p.char == charNewline {
			p.firstSpace = charNewline
		}
		// lookahead by 1, check for space
		if p.pos+1 < p.len && unicode.IsSpace(p.input[p.pos+1]) {
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
	repeat:
		switch p.char {
		case eof, terminator:
			return
		case charBackslash:
			if p.pos+1 < p.len {
				escapeChar := p.input[p.pos+1]
				p.ignoreChar = true
				p.next()
				p.ignoreChar = false
				p.char = escapeChar
				p.next()
				goto repeat
			}
		}
		if strings.ContainsRune(destinationOptions, p.char) {
			return
		}
		p.next()
	}
}

func (p *Parser) nextOverKeyword() bool {
	if p.char == eof || p.char == terminator {
		return false
	}
	for {
		switch p.char {
		case eof, terminator:
			return false
		}
		if !unicode.IsLetter(p.char) && !unicode.IsNumber(p.char) {
			return true
		}
		p.next()
	}
}

func (p *Parser) skipWhitespace() {
	if unicode.IsSpace(rune(p.char)) {
		p.next()
	}
}

// func (p *Parser) debug() {
// 	for {
// 		fmt.Printf("%q,%d,%q,%d\n", p.frame, p.pos, string(p.char), p.char)
// 		if p.char == eof {
// 			os.Exit(1)
// 		}
// 		p.next()
// 	}
// }

func trim(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
}
