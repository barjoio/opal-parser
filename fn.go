package opalparser2

type parseFn func(*Parser) parseFn

// parseBegin begins parsing the next major element
// this will be either Paragraph or BlockTag
func parseBegin(p *Parser) parseFn {
	p.next()
	p.skipWhitespace(true)

	switch p.char {
	case eof:
		p.addToParent()
		return nil
	case charFullstop:
		return nil
	default:
		p.createParentNode(nodeParagraph)
		return parseParagraph
	}
}

// parseParagraph parses Paragraph elements
// can contain: Text, InlineTag
func parseParagraph(p *Parser) parseFn {
	p.nextUntil(nodeParagraphText, "`", "")

	switch p.char {
	case eof:
		p.addToParent()
		return nil
	case terminator:
		p.addToParent()
		return parseBegin
	}

	p.next() // skip over grave
	p.flattenFrame()

	p.pushParseFn()
	p.createParentNode(nodeInlineTag)
	return parseInlineTag
}

// parseInlineTag parses InlineTag elements
// can contain: InlineTagName, InlineTagText
func parseInlineTag(p *Parser) parseFn {

	// leading whitespace

	p.skipWhitespace(false)
	switch p.char {
	case eof:
		p.addError(errUnexpectedEOF)
		p.addToParent()
		return p.popParseFn()
	case terminator:
		p.addError(errUnexpectedTerm)
		p.addToParent()
		p.next()
		p.flattenFrame()
		return p.popParseFn()
	}

	// get tag name

	p.nextUntil(nodeInlineTagName, charsWhitespace, "`")
	switch p.char {
	case eof:
		p.addError(errUnexpectedEOF)
		p.addToParent()
		return p.popParseFn()
	case terminator:
		p.addError(errUnexpectedTerm)
		p.addToParent()
		p.next()
		p.flattenFrame()
		return p.popParseFn()
	}

	// get tag content

	p.nextUntil(nodeInlineTagText, "`", "")
	switch p.char {
	case eof:
		p.addError(errUnexpectedEOF)
		p.addToParent()
		return p.popParseFn()
	case terminator:
		p.addError(errUnexpectedTerm)
		p.addToParent()
		p.next()
		p.flattenFrame()
		return p.popParseFn()
	}

	p.next() // skip over grave
	p.flattenFrame()

	p.addToParent()
	return p.popParseFn()
}
