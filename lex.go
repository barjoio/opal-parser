package opalparser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// token represents a token or text string returned from the scanner.
type token struct {
	typ  tokenType // The type of this token.
	pos  int       // The starting position, in bytes, of this token in the input string.
	val  string    // The value of this token.
	line int       // The line number at the start of this token.
}

// print format for tokens
func (t token) String() string {
	fmt.Printf("[%d,%d]: ", t.line, t.pos+1)
	if t.typ == tokenError {
		fmt.Printf("Erroneous char: ")
	}
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenNewline:
		return "New line"
	case tokenSpace:
		return "Whitespace"
	case tokenTerminator:
		if t.val == string(charSemicolon) {
			return "Terminator: ';'"
		}
		return "Terminator: Doubleline"
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

// tokenType identifies the type of lex items.
type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenSpace

	tokenFullstop
	tokenColon
	tokenSlash
	tokenNewline
	tokenSemicolon
	tokenGrave
	tokenGateOpen
	tokenGateClose
	tokenTerminator

	tokenElement
	tokenAttr
	tokenInlineBlock
	tokenStdBlock
	tokenExtBlock

	tokenParagraphText
	tokenInlineElement
	tokenInlineElementContent
)

const (
	eof           rune = -1
	charFullstop       = '.'
	charColon          = ':'
	charSlash          = '/'
	charNewline        = '\n'
	charSemicolon      = ';'
	charGrave          = '`'
	charGateOpen       = '-'
	charGateClose      = "----"
)

// lexFn represents the state of the scanner as a function that returns the next state.
type lexFn func(*lexer) lexFn

type lexer struct {
	name      string     // the name of the input; used only for error reports
	input     string     // the string being scanned
	pos       int        // current position in the input
	start     int        // start position of this item
	width     int        // width of last rune read from input
	tokens    chan token // channel of scanned tokens
	line      int        // 1+number of newlines seen
	startLine int        // start line of this item
	len       int
	termLen   int
	cur       rune
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= l.len {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	l.cur = r
	return r
}

func (l *lexer) nextN(n int) {
	if n > 0 {
		for i := 0; i < n; i++ {
			l.next()
		}
	}
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
	l.cur = rune(l.input[l.pos])
}

func (l *lexer) backupN(n int) {
	if n > 0 {
		for i := 0; i < n; i++ {
			l.backup()
		}
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	if l.pos > l.start || t == tokenEOF {
		l.tokens <- token{t, l.start, l.input[l.start:l.pos], l.startLine}
	}
	l.start = l.pos
	l.startLine = l.line
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.startLine}
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() token {
	return <-l.tokens
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:      name,
		input:     input + string(charSemicolon),
		tokens:    make(chan token),
		line:      1,
		startLine: 1,
		len:       len(input),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexBegin; state != nil; {
		state = state(l)
	}

	l.emit(tokenEOF)
	close(l.tokens)
}

//
//
//

func (l *lexer) posToEnd() string {
	return l.input[l.pos:]
}

func (l *lexer) isEOF() bool {
	return l.pos >= l.len
}

func (l *lexer) accomNewLines() {
	l.line += strings.Count(l.input[l.start:l.pos], string(charNewline))
}

func (l *lexer) skipSpace() {
	for i := 0; ; i++ {
		r := l.next()
		if !isSpace(r) {
			if r != eof {
				l.backup()
			}
			if i > 0 {
				l.emit(tokenSpace)
			}
			break
		}
	}
}

// checkForTerminator checks to see if the current rune, or runes
// thereafter, represent a termination. This is either a semicolon
// or a group of whitespace containing 2 newlines
func (l *lexer) checkForTerminator() rune {
	r := l.cur
	if r != charSemicolon && r != charNewline {
		return -1
	}
	if r == charSemicolon {
		l.termLen = 1
		return r
	}
	diff := len(l.posToEnd()) - len(strings.TrimSpace(l.posToEnd()))
	if diff >= 1 {
		diffText := l.input[l.pos : l.pos+diff]
		if strings.Count(diffText, "\n") >= 1 {
			l.termLen = diff + 1
			return r
		}
	}
	return -1
}

func (l *lexer) checkForSpace() rune {
	if isSpace(l.cur) {
		return l.cur
	}
	return -1
}

//
//	state funcs
//

func lexBegin(l *lexer) lexFn {
	l.skipSpace()

	switch l.next() {
	case eof:
		return nil
	case charFullstop:
		l.emit(tokenFullstop)
		return lexElement
	}
	return lexParagraph
}

func lexElement(l *lexer) lexFn {
	for {
		switch l.next() {
		case eof:
			l.emit(tokenElement)
			return nil
		case charColon:
			l.backup()
			l.emit(tokenElement)
			l.next()
			l.emit(tokenColon)
			return lexInlineBlock
		case charSlash:
			return nil
		case l.checkForTerminator():
			l.backup()
			l.emit(tokenElement)
			l.nextN(l.termLen)
			l.emit(tokenTerminator)
			return lexBegin
		case charNewline:
			l.backup()
			l.emit(tokenElement)
			l.next()
			l.emit(tokenNewline)
			return nil
		}
	}
}

func lexInlineBlock(l *lexer) lexFn {
	for {
		switch l.next() {
		case eof:
			l.emit(tokenInlineBlock)
			return nil
		case l.checkForTerminator():
			l.backup()
			l.emit(tokenInlineBlock)
			l.nextN(l.termLen)
			l.emit(tokenTerminator)
			return lexBegin
		case charNewline:
			l.backup()
			l.emit(tokenElement)
			l.next()
			l.emit(tokenNewline)
		}
	}
}

func lexParagraph(l *lexer) lexFn {
	for {
		switch l.next() {
		case eof:
			l.emit(tokenParagraphText)
			return nil
		case charGrave:
			l.backup()
			l.emit(tokenParagraphText)
			l.next()
			l.emit(tokenGrave)
			return lexInlineElement
		case l.checkForTerminator():
			l.backup()
			l.emit(tokenParagraphText)
			l.nextN(l.termLen)
			l.emit(tokenTerminator)
			return lexBegin
		case charNewline:
			l.backup()
			l.emit(tokenParagraphText)
			l.next()
			l.emit(tokenNewline)
		}
	}
}

func lexInlineElement(l *lexer) lexFn {
	for {
		switch l.next() {
		case eof:
			l.emit(tokenInlineElement)
			return nil
		case charGrave:
			l.backup()
			l.emit(tokenInlineElement)
			l.next()
			l.emit(tokenGrave)
			return lexParagraph
		case l.checkForTerminator():
			l.backup()
			l.emit(tokenInlineElement)
			l.nextN(l.termLen)
			l.emit(tokenTerminator)
			return lexBegin
		case l.checkForSpace():
			l.backup()
			l.emit(tokenInlineElement)
			l.next()
			l.emit(tokenSpace)
			return lexInlineElementContent
		}
	}
}

func lexInlineElementContent(l *lexer) lexFn {
	for {
		switch l.next() {
		case eof:
			l.emit(tokenInlineElementContent)
			return nil
		case charGrave:
			l.backup()
			l.emit(tokenInlineElementContent)
			l.next()
			l.emit(tokenGrave)
			return lexParagraph
		}
	}
}
