package opalparser

import (
	"strings"
	"unicode"
)

type parseFn func(*Parser) parseFn

// parseBegin begins parsing the next major element
// this will be either Paragraph or BlockTag
func parseBegin(p *Parser) parseFn {
repeat:
	switch p.char {
	case eof:
		return nil
	case terminator:
		p.nextFlat()
		goto repeat
	case charFullstop:
		return parseBlockTag
	}
	if unicode.IsSpace(p.char) {
		p.nextFlat()
		goto repeat
	}
	return parseParagraph
}

// parseBlockTag parses BlockTag elements
// can contain: BlockTagName, BlockTagText
func parseBlockTag(p *Parser) parseFn {
	p.createNode(nodeBlockTag)

	p.next() // skip over fullstop
	p.flattenFrame()

	p.skipWhitespace()

	// get block tag name
	p.nextOverKeyword()
	p.determineNodeType()
	switch p.char {
	case eof:
		p.addErrorUnexpected()
		return nil
	case terminator:
		p.addToParent()
		return parseBegin
	}
	if len(p.frame) == 0 {
		p.addErrorUnexpected()
		return parseBegin
	}

	// set heading level
	switch string(p.frame) {
	case "1", "2", "3", "4", "5", "6":
		p.currentNode().Level = string(p.frame)
	}

	p.skipWhitespace()

	switch p.char {
	case eof:
		p.addToParent()
		return nil
	case terminator:
		p.addToParent()
		return parseBegin
	case charColon:
		p.nextFlat() // skip over colon
		parseText(p, "\n", nil)
		p.addToParent()
		return parseBegin
	case charSlash:
		return parseBlockTagAttrs
	case charNewline:
		p.nextFlat()
		parseStdBlock(p)
		p.addToParent()
		return parseBegin
	}
	p.addErrorUnexpected()
	return parseBegin
}

// parseBlockTagAttrs parses BlockTag attributes
// can contain: BlockTagAttr
func parseBlockTagAttrs(p *Parser) parseFn {
repeat:
	p.nextFlat()
	p.nextOverKeyword()
	p.appendAttr()
	switch p.char {
	case eof:
		p.addToParent()
		return nil
	case terminator:
		p.addToParent()
		p.nextFlat()
		return parseBegin
	case charSlash:
		goto repeat
	case charNewline:
		p.nextFlat()
		parseStdBlock(p)
		p.addToParent()
		return parseBegin
	}
	p.addErrorUnexpected()
	return parseBegin
}

func parseStdBlock(p *Parser) {
	switch p.nodeType {
	case nodeList:
		p.createNode(nodeListItem)
		parseText(p, "-", func() {
			p.addChild(nodeText, true, false)
			p.addPopulatedToParent()
			p.createNode(nodeListItem)
		})
		p.addPopulatedToParent()
	case nodeTable:
		p.createNode(nodeTableRow)
		p.createNode(nodeTableData)
		parseText(p, "|\n", func() {
			p.addChild(nodeText, true, false)
			// table data
			p.addPopulatedToParent()
			// table row
			if p.char == charNewline {
				p.addToParent()
				p.createNode(nodeTableRow)
			}
			// new table data
			p.createNode(nodeTableData)
		})
		// tableData
		p.addPopulatedToParent()
		// tableRow
		p.addPopulatedToParent()
	default:
		parseText(p, "", nil)
	}
}

// parseParagraph parses paragraphs
// can contain: Text, InlineTag
func parseParagraph(p *Parser) parseFn {
	p.createNode(nodeParagraph)
	parseText(p, "", nil)
	p.addToParent()
	return parseBegin
}

// parseText parses text content
// can contain: Text, InlineTag
func parseText(p *Parser, splitOn string, callback func()) {
repeat:
	p.nextUntil("`" + string(splitOn))
	switch p.char {
	case eof, terminator:
		if callback != nil {
			callback()
		} else {
			p.addChild(nodeText, true, true)
		}
		return
	}
	if strings.ContainsRune(splitOn, p.char) {
		if callback != nil {
			callback()
		} else {
			p.addChild(nodeText, true, true)
		}
		p.nextFlat()
		goto repeat
	}
	p.addChild(nodeText, true, true)
	parseInlineTag(p)
	goto repeat
}

// parseInlineTag parses InlineTag elements
// can contain: InlineTagName, InlineTagText
func parseInlineTag(p *Parser) {
	p.createNode(nodeInlineTag)

	p.nextFlat()       // skip over grave
	p.skipWhitespace() // skip leading space
	p.flattenFrame()

	// get inline tag name
	p.nextOverKeyword()
	p.determineNodeType()
	switch p.char {
	case eof, terminator:
		p.addErrorUnexpected()
		return
	}
	if !unicode.IsSpace(p.char) {
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

	// check for empty tag text
	if len(p.frame) == 0 {
		p.addErrorUnexpected()
		return
	}

	val := string(p.frame)

	switch p.nodeType {
	case nodeHyperlink:
		parseHyperlink(p, val)
	default:
		p.currentNode().Value = val
	}

	p.addToParent()
	p.next()
	p.flattenFrame()
	return
}

func parseHyperlink(p *Parser, s string) {
	li := strings.LastIndexAny(s, charsWhitespace)
	if li == -1 {
		p.addError(errNoURL)
		p.currentNode().DisplayText = s
		return
	}
	url := string(p.frame[li+1:])
	displayText := string(p.frame[:li])
	if displayText == "_" {
		displayText = url
	}
	p.currentNode().DisplayText = displayText
	p.currentNode().URL = url
}
