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
	switch {
	case t.typ == tokenEOF:
		return "EOF"
	case t.typ == tokenNewline:
		return "New line"
	case t.typ == tokenSpace:
		return "Whitespace"
	case t.typ == tokenTerminator:
		return "Terminator"
	case len(t.val) > 10:
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
	charFullstop  = "."
	charColon     = ":"
	charSlash     = "/"
	charNewline   = "\n"
	charSemicolon = ";"
	charGrave     = "`"
	charGateOpen  = "-"
	charGateClose = "----"
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
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

func (l *lexer) nextN(n int) {
	if n > 0 {
		for i := 0; i < n; i++ {
			l.next()
		}
	}
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backupN(n int) {
	if n > 0 {
		for i := 0; i < n; i++ {
			l.backup()
		}
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos], l.startLine}
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
func (l *lexer) errorf(format string, args ...interface{}) lexFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.startLine}
	return nil
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
		input:     input + charSemicolon,
		tokens:    make(chan token),
		line:      1,
		startLine: 1,
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexBegin; state != nil && l.pos < len(l.input)-1; {
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

func (l *lexer) accomNewLines() {
	l.line += strings.Count(l.input[l.start:l.pos], charNewline)
}

func (l *lexer) skipSpace() {
	l.start = l.pos
	for isSpace(l.next()) {
	}
	l.backup()
	if l.pos-l.start > 0 {
		l.emit(tokenSpace)
	}
}

func (l *lexer) checkForTerminator(callback func()) bool {
	if strings.HasPrefix(l.posToEnd(), charSemicolon) {
		callback()
		l.next()
		l.nextN(len(l.posToEnd()) - len(strings.TrimSpace(l.posToEnd())))
		l.emit(tokenTerminator)
		return true
	} else if strings.HasPrefix(l.posToEnd(), charNewline) {
		diff := len(l.posToEnd()) - len(strings.TrimSpace(l.posToEnd()))
		if diff >= 2 && strings.Count(l.input[l.pos:l.pos+diff], charNewline) >= 2 {
			callback()
			l.nextN(diff)
			l.emit(tokenTerminator)
			return true
		}
	}
	return false
}

//
//	state funcs
//

func lexBegin(l *lexer) lexFn {
	l.skipSpace()

	if strings.HasPrefix(l.posToEnd(), charFullstop) {
		return lexFullstop
	}
	return lexParagraph
}

func lexFullstop(l *lexer) lexFn {
	l.next()
	l.emit(tokenFullstop)
	return lexElement
}

func lexElement(l *lexer) lexFn {
	for {
		if strings.HasPrefix(l.posToEnd(), charColon) {
			l.emit(tokenElement)
			l.next()
			l.emit(tokenColon)
			l.skipSpace()
			return lexInlineBlock
		} else if strings.HasPrefix(l.posToEnd(), charSlash) {
			l.emit(tokenElement)
			return lexBegin
		} else if l.checkForTerminator(func() { l.emit(tokenElement) }) {
			return lexBegin
		} else if strings.HasPrefix(l.posToEnd(), charNewline) {
			l.emit(tokenElement)
			l.next()
			l.emit(tokenNewline)
			l.next()
			return lexBegin
		}

		l.next()
	}
}

func lexInlineBlock(l *lexer) lexFn {
	if l.checkForTerminator(func() { l.emit(tokenInlineBlock) }) {
		return lexBegin
	}
	l.next()
	return lexInlineBlock
}

func lexParagraph(l *lexer) lexFn {
	if strings.HasPrefix(l.posToEnd(), charGrave) {
		l.accomNewLines()
		l.emit(tokenParagraphText)
		l.next()
		l.emit(tokenGrave)
		return lexInlineElement
	} else if l.checkForTerminator(func() { l.accomNewLines(); l.emit(tokenParagraphText) }) {
		return lexBegin
	}
	l.next()
	return lexParagraph
}

func lexInlineElement(l *lexer) lexFn {
	if strings.HasPrefix(l.posToEnd(), charGrave) {
		l.accomNewLines()
		l.emit(tokenInlineElement)
		l.next()
		l.emit(tokenGrave)
		return lexParagraph
	}
	l.next()
	return lexInlineElement
}
