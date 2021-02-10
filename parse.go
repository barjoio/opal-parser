// https://blog.gopheracademy.com/advent-2014/parsers-lexers/
// https://medium.com/@bradford_hamilton/building-a-json-parser-and-query-tool-with-go-8790beee239a
// http://blog.leahhanson.us/post/recursecenter2016/recipeparser.html

package opalparser

import (
	"fmt"
	"strconv"
)

// Parser ...
type Parser struct {
	l            *lexer
	currentToken token
	currentItem  Item
	prevToken    token
	items        []Item
}

// New ...
func New() *Parser {
	return &Parser{
		prevToken: token{typ: -1},
	}
}

// Item ...
type Item struct {
	Typ string
	Val string
}

// Parse ...
func (p *Parser) Parse(input string) []Item {
	p.l = lex("", input)

	var t token
	for t.typ != tokenEOF {
		t = p.l.nextItem()
		fmt.Println(t)
	}
	// close(p.l.tokens)

	// for state := parseBegin; state != nil; {
	// 	state = state(p)
	// }

	return p.items
}

func (p *Parser) addItem() {
	p.items = append(p.items, p.currentItem)
	p.currentItem = Item{}
}

func (p *Parser) skipSpace() token {
	t := p.l.nextItem()
	if t.typ == tokenSpace {
		return p.skipSpace()
	}
	return t
}

type parseFn func(p *Parser) parseFn

func parseBegin(p *Parser) parseFn {
	t := p.skipSpace()
	switch t.typ {
	case tokenFullstop:
		return parseElement
	default:
		p.currentToken = t
		return parseParagraph
	}
}

func parseElement(p *Parser) parseFn {
	t := p.l.nextItem()
	switch t.typ {
	case tokenElement:
		_, err := strconv.Atoi(t.val)
		if err == nil || t.val == "image" {
			p.currentItem.Typ = t.val
			t = p.l.nextItem()
			switch t.typ {
			case tokenColon:
				t = p.skipSpace()
				switch t.typ {
				case tokenInlineBlock:
					p.currentItem.Val = t.val
					p.addItem()
					t = p.l.nextItem()
					switch t.typ {
					case tokenTerminator:
						p.currentItem.Typ = "Terminator"
						p.currentItem.Val = ""
						p.addItem()
						return parseBegin
					}
				}
			}
		}
	}
	return nil
}

func parseParagraph(p *Parser) parseFn {
	switch p.currentToken.typ {
	case tokenParagraphText:
		p.currentItem.Typ = "Paragraph"
		p.currentItem.Val = p.currentToken.val
		p.addItem()
		p.currentToken = p.l.nextItem()
		return parseParagraph
	case tokenGrave:
		t := p.l.nextItem()
		switch t.typ {
		case tokenInlineElement:
			var addEl bool
			if len(t.val) >= 3 {
				if (string(t.val[0]) == "b" || string(t.val[0]) == "i") && string(t.val[1]) == " " {
					addEl = true
					p.currentItem.Typ = string(t.val[0])
					p.currentItem.Val = string(t.val[2:])
				}
			} else {
				return nil
			}
			t = p.l.nextItem()
			switch t.typ {
			case tokenGrave:
				if addEl {
					p.addItem()
				}
				p.currentToken = p.l.nextItem()
				return parseParagraph
			}
		}
	case tokenTerminator:
		p.currentItem.Typ = "Terminator"
		p.currentItem.Val = ""
		p.addItem()
		return parseBegin
	}
	return nil
}
