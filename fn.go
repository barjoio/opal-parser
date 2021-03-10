package opalparser2

import (
	"strings"
)

type parseFn func(*Parser) parseFn

// parseBegin begins parsing the next major element
// this will be either Paragraph or BlockTag
func parseBegin(p *Parser) parseFn {
	p.next()

	switch p.char {
	case eof:
		return nil
	case charFullstop:
		p.createNode(nodeBlockTag)
		return parseBlockTag
	}
	return parseParagraph
}

// parseBlockTag parses BlockTag elements
// can contain: BlockTagName, BlockTagText
func parseBlockTag(p *Parser) parseFn {
	// skip over fullstop
	p.next()
	p.flattenFrame()

	// leading whitespace

	p.skipWhitespace()
	switch p.char {
	case eof:
		p.addError(errUnexpectedEOF)
		p.addToParent()
		return nil
	case terminator:
		p.addError(errUnexpectedTerm)
		p.addToParent()
		p.next()
		p.flattenFrame()
		return parseBegin
	}

	// get block tag name

	p.nextUntil(":/\n")
	switch p.char {
	case eof:
		p.addError(errUnexpectedEOF)
		p.addToParent()
		return nil
	case terminator:
		p.addToParent()
		p.next()
		p.flattenFrame()
		return parseBegin
	case charColon:
		p.next() // skip over colon
		p.flattenFrame()
		p.createNode(nodeBlockTagLine)
		return parseBlockTagText
	case charSlash:
		return parseBlockTagAttrs
	case charNewline:
		p.next()
		p.flattenFrame()
		p.createNode(nodeBlockTagLine)
		return parseBlockTagText
	}
	return parseBegin
}

// get single line tag text
func parseBlockTagText(p *Parser) parseFn {
	for {
		p.nextUntil("`\n")
		switch p.char {
		case eof:
			if p.numChildren() == 0 {
				p.popNode()
			} else {
				p.addToParent() // line
			}
			p.addToParent() // add nodeBlockTag
			return nil
		case terminator:
			if p.numChildren() == 0 {
				p.popNode()
			} else {
				p.addToParent() // line
			}
			p.addToParent() // add nodeBlockTag
			p.next()
			p.flattenFrame()
			return parseBegin
		case charGrave:
			p.next() // skip over grave
			p.flattenFrame()
			p.pushParseFn()
			p.createNode(nodeInlineTag)
		case charNewline:
			p.addToParent() // line
			p.next()
			p.flattenFrame()
			p.createNode(nodeBlockTagLine)
		}
	}
}

// parseBlockTagAttrs parses BlockTag attributes
// can contain: BlockTagAttr
func parseBlockTagAttrs(p *Parser) parseFn {
	for {
		p.next()
		p.flattenFrame()
		p.nextUntil("/\n")
		switch p.char {
		case eof:
			p.addError(errUnexpectedEOF)
			p.addToParent()
			return nil
		case terminator:
			p.addToParent()
			p.next()
			p.flattenFrame()
			return parseBegin
		case charSlash:
			continue
		case charNewline:
			p.next()
			p.flattenFrame()
			p.createNode(nodeBlockTagLine)
			return parseBlockTagText
		}
	}
}

// parseParagraph parses Paragraph elements
// can contain: Text, InlineTag
func parseParagraph(p *Parser) parseFn {
	p.createNode(nodeParagraph)

repeat:
	p.nextUntil("`")
	p.addChild(nodeText)
	switch p.char {
	case eof:
		p.addToParent()
		return nil
	case terminator:
		p.nextFlat()
		p.addToParent()
		return parseBegin
	}
	parseInlineTag(p)
	goto repeat
}

// parseInlineTag parses InlineTag elements
// can contain: InlineTagName, InlineTagText
func parseInlineTag(p *Parser) {
	p.createNode(nodeInlineTag)

	p.nextFlat()
	p.skipWhitespace()
	p.flattenFrame()

	// get inline tag name
	p.nextUntil(charsWhitespace + "`")
	p.determineNodeType()
	switch p.char {
	case eof, terminator, charGrave:
		p.addErrorUnexpected()
		return
	}

	// skip over whitespace separator
	p.nextFlat()

	// get tag text
	p.nextUntil("`")
	switch p.char {
	case eof, terminator:
		p.addErrorUnexpected()
		return
	}

	p.trimFrame()

	if p.frame == "" {
		p.addErrorUnexpected()
		return
	}

	switch p.nodeType {
	case nodeHyperlink:
		parseHyperlink(p)
	default:
		p.currentNode().Value = p.frame
	}

	p.addToParent()
	p.next()
	p.flattenFrame()
	return
}

func parseHyperlink(p *Parser) {
	li := strings.LastIndexAny(p.frame, charsWhitespace)
	if li == -1 {
		p.addError(errNoURL)
		p.currentNode().DisplayText = p.frame
		return
	}
	url := p.frame[li+1:]
	displayText := p.frame[:li]
	if displayText == "_" {
		displayText = url
	}
	p.currentNode().DisplayText = displayText
	p.currentNode().URL = url
}
