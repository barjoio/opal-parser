package opalparser

import "strings"

// Node is a grammatically defined element in the Opal language
// these are used to construct the abstract syntax tree
type Node struct {
	Typ         nodeType  `json:"type,omitempty"`
	Errors      []errType `json:"errors,omitempty"`
	Value       string    `json:"value,omitempty"`
	Attrs       []string  `json:"attrs,omitempty"`
	DisplayText string    `json:"displayText,omitempty"`
	URL         string    `json:"url,omitempty"`
	Level       string    `json:"level,omitempty"`
	Ln          int       `json:"line,omitempty"`
	Col         int       `json:"column,omitempty"`
	Children    []*Node   `json:"children,omitempty"`
}

type nodeType int

// list of nodes that can be parsed
const (
	nodeEOF nodeType = iota
	nodeInvalidTag
	nodeWhitespace
	nodeRoot
	nodeText
	nodeListItem
	nodeTagName
	nodeBlockTag
	nodeBlockTagLine
	nodeAttr
	nodeParagraph
	nodeList
	nodeTable
	nodeTitle
	nodeToC
	nodeHeading
	nodeInlineTag
	nodeHyperlink
	nodeBoldText
	nodeItalicText
	nodeUnderlineText
	nodeBoldItalic
	nodeBoldUnderline
	nodeItalicUnderline
	nodeCode
	nodeTableRow
	nodeTableData
)

// makeNode returns a new node
// parent nodes have only a type and list of children
func (p *Parser) makeNode(n nodeType, hasVal, hasLineInfo bool) *Node {
	val := ""
	var startLn, startCol int
	if hasVal {
		val = trim(string(p.frame))
	}
	if hasLineInfo {
		startLn = p.startLn
		startCol = p.startCol
	}
	return &Node{Typ: n, Value: val, Ln: startLn, Col: startCol}
}

// createNode appends a new parent node to the node stack
func (p *Parser) createNode(n nodeType) {
	p.nodeStack = append(p.nodeStack, p.makeNode(n, false, true))
	p.nodeType = n
	p.nodesParsedLinear = append(p.nodesParsedLinear, n)
}

// currentNode returns the topmost node from the node stack
func (p *Parser) currentNode() *Node {
	if len(p.nodeStack) > 0 {
		return p.nodeStack[len(p.nodeStack)-1]
	}
	return nil
}

// lastNode returns get the most recently added node
func (p *Parser) lastNode() *Node {
	topNodeChildren := p.currentNode().Children
	topNodeChildrenLen := len(topNodeChildren)
	if topNodeChildrenLen == 0 {
		return nil
	}
	return topNodeChildren[topNodeChildrenLen-1]
}

// addChild appends a new child node to the topmost node from the node stack
func (p *Parser) addChild(n nodeType, merge, hasLineInfo bool) {
	// if the last node is of the same type, add the frame content
	// to the end of the last node
	val := trim(string(p.frame))
	if val != "" {
		lastNode := p.lastNode()
		if merge && lastNode != nil && lastNode.Typ == n && n != nodeListItem {
			p.lastNode().Value += " " + val
		} else {
			topNode := p.currentNode()
			topNode.Children = append(topNode.Children, p.makeNode(n, true, hasLineInfo))
		}
	}
	p.flattenFrame()
	p.nodesParsedLinear = append(p.nodesParsedLinear, n)
}

// popNode removes the topmost node from the node stack
func (p *Parser) popNode() {
	p.nodeStack = p.nodeStack[:len(p.nodeStack)-1]
}

// addToParent adds the topmost node from the node stack to the abstract syntax tree and removes it from the node stack
func (p *Parser) addToParent() {
	topNode := p.currentNode()
	lenNodeStack := len(p.nodeStack)

	if lenNodeStack == 0 {
		return
	}

	if lenNodeStack == 1 {
		p.tree = append(p.tree, topNode)
		p.popNode()
		return
	}

	secondToTopNode := p.nodeStack[lenNodeStack-2]
	secondToTopNode.Children = append(secondToTopNode.Children, topNode)
	p.popNode()
}

func (p *Parser) addPopulatedToParent() {
	if len(p.currentNode().Children) > 0 {
		p.addToParent()
	} else {
		p.popNode()
	}
}

func (p *Parser) determineNodeType() {
	if len(p.frame) == 0 {
		return
	}
	var t nodeType
	switch strings.ToLower(string(p.frame)) {
	case "1", "2", "3", "4", "5", "6":
		t = nodeHeading
	case "b":
		t = nodeBoldText
	case "bi", "ib":
		t = nodeBoldItalic
	case "bu", "ub":
		t = nodeBoldUnderline
	case "c":
		t = nodeCode
	case "i":
		t = nodeItalicText
	case "iu", "ui":
		t = nodeItalicUnderline
	case "l":
		t = nodeHyperlink
	case "list":
		t = nodeList
	case "table":
		t = nodeTable
	case "toc":
		t = nodeToC
	case "title":
		t = nodeTitle
	case "u":
		t = nodeUnderlineText
	default:
		t = nodeInvalidTag
		p.addError(errInvalidTagName)
	}
	p.currentNode().Typ = t
	p.nodeType = t
}

func (p *Parser) appendAttr() {
	if len(p.frame) > 0 {
		attrs := &p.currentNode().Attrs
		*attrs = append(*attrs, string(p.frame))
	}
}
