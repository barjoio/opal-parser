package opalparser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Parser is used for processing Opal markup
type Parser struct {
	contextStack []*context
	contextMap   map[string]func(i int, r rune)

	lineNum int
	linePos int

	markup string
}

// Element represents an element in the markup
type Element struct {
	Name     string     `json:"name,omitempty"`
	LineNum  int        `json:"line_num,omitempty"`
	LinePos  int        `json:"line_pos,omitempty"`
	Value    string     `json:"value,omitempty"`
	Children []*Element `json:"children,omitempty"`
}

// context is used in the contextStack to modulate character processing
type context struct {
	name           string
	currentElement *Element
	elementPos     int
	textBuffer     string
}

// New returns a new Parser
func New() (p *Parser) {
	return &Parser{
		contextStack: []*context{
			{name: "start", currentElement: &Element{Name: "opal"}},
		},
		contextMap: map[string]func(i int, r rune){
			"start": func(i int, r rune) {
				if r == '!' {
					newElement := p.createElement("tag", "")
					p.addChildElement(newElement)
					p.pushContext(p.createContext("tagName", newElement))
				} else if !isSpace(r) {
					newElement := p.createElement("paragraph", "")
					p.addChildElement(newElement)
					p.pushContext(p.createContext("text", newElement))
					p.appendTextBuffer(r)
				}
			},
			"tagName": func(i int, r rune) {
				if r == '\n' {
					p.addChildElement(p.createElement("tagName", p.getTextBuffer()))
					p.popContext()
				} else {
					p.appendTextBuffer(r)
				}
			},
			"text": func(i int, r rune) {
				if r == '\n' {
					if p.getTextBuffer() != "" {
						p.addChildElement(p.createElement("text", p.getTextBuffer()))
					}
					p.popContext()
				} else if r == '(' && isLetter(rune(p.markup[i-1])) {

					// lookbehind for tagName
					var tagName string
					for {
						i--
						r := rune(p.markup[i])
						if !isLetter(r) {
							break
						} else {
							tagName = string(r) + tagName
						}
					}

					// remove tagName from end of textBuffer
					textBuffer := p.getTextBuffer()
					textBuffer = textBuffer[:len(textBuffer)-len(tagName)]

					lenTagName := len(tagName)

					// remove fullstop from end of textBuffer if used in concatenation
					if len(textBuffer) > 0 && textBuffer[len(textBuffer)-1] == '.' {
						textBuffer = textBuffer[:len(textBuffer)-1]
						lenTagName++
					}

					// add text element
					currentlinePos := p.linePos
					p.linePos -= lenTagName
					if textBuffer != "" {
						p.addChildElement(p.createElement("text", textBuffer))
					}

					// add inlineTag element
					newElement := p.createElement("inlineTag", "")
					p.addChildElement(newElement)
					p.pushContext(p.createContext("inlineTagAttrs", newElement))

					// add inlineTagName element to inlineTag
					p.linePos = currentlinePos
					p.addChildElement(p.createElement("inlineTagName", tagName))
				} else {
					p.appendTextBuffer(r)
				}
			},
			"inlineTagAttrs": func(i int, r rune) {
				if r == ')' {
					p.addChildElement(p.createElement("inlineTagAttr", strings.TrimSpace(p.getTextBuffer())))
					p.popContext()
				} else if r == ',' {
					p.addChildElement(p.createElement("inlineTagAttr", strings.TrimSpace(p.getTextBuffer())))
				} else {
					p.appendTextBuffer(r)
				}
			},
		},
	}
}

// Parse processes a string of markup and returns an object model of the markup
func (p *Parser) Parse(markup string) *Element {
	p.markup = "\n" + markup + "\n"
	p.lineNum = 0
	p.linePos = 0

	for i, c := range p.markup {
		r := rune(c)

		// update line number and position
		p.linePos++
		if r == '\n' {
			p.lineNum++
			p.linePos = 0
		}

		// execute the current context's contextAction on the current character (try saying that twice as fast ;-))
		if contextAction, exists := p.contextMap[p.getContext().name]; exists {
			contextAction(i, r)
		} else {
			fmt.Println("Unknown context")
		}
	}

	return p.contextStack[0].currentElement
}

// ParseFile accepts a filepath to parse markup from
func (p *Parser) ParseFile(filepath string) *Element {
	b, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	return p.Parse(string(b))
}

// context operations

func (p *Parser) pushContext(c *context) {
	c.elementPos = p.linePos
	p.contextStack = append(p.contextStack, c)
}

func (p *Parser) popContext() {
	l := len(p.contextStack)
	if l > 0 {
		p.contextStack = p.contextStack[:l-1]
	}
}

func (p *Parser) getContext() *context {
	l := len(p.contextStack)
	if l > 0 {
		return p.contextStack[l-1]
	}
	return nil
}

func (p *Parser) createContext(name string, currentElement *Element) *context {
	return &context{name: name, currentElement: currentElement}
}

// element operations

func (p *Parser) getCurrentElement() *Element {
	return p.getContext().currentElement
}

func (p *Parser) createElement(name, value string) *Element {
	return &Element{
		Name:    name,
		LineNum: p.lineNum,
		LinePos: p.linePos - len(value),
		Value:   value,
	}
}

func (p *Parser) addChildElement(newElement *Element) {
	currentElement := p.getContext().currentElement
	currentElement.Children = append(currentElement.Children, newElement)
	p.setTextBuffer("")
}

// textBuffer operations

func (p *Parser) setTextBuffer(t string) {
	p.getContext().textBuffer = t
}

func (p *Parser) appendTextBuffer(r rune) {
	p.getContext().textBuffer += string(r)
}

func (p *Parser) getTextBuffer() string {
	return p.getContext().textBuffer
}
