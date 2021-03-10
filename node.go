package opalparser2

import (
	"strings"
)

// Node is a grammatically defined element in the Opal language
// these are used to construct the abstract syntax tree
type Node struct {
	Typ         nodeType  `json:"type,omitempty"`
	Errors      []errType `json:"errors,omitempty"`
	Value       string    `json:"value,omitempty"`
	DisplayText string    `json:"displayText,omitempty"`
	URL         string    `json:"url,omitempty"`
	Ln          int       `json:"line,omitempty"`
	Col         int       `json:"column,omitempty"`
	Children    []*Node   `json:"children,omitempty"`
}

// nodeData is an object describing data relating to a node
type nodeData struct {
	typ   string
	value string
}

type nodeType string

// list of nodes that can be parsed
const (
	nodeEOF          nodeType = "EOF"
	nodeInvalidTag   nodeType = "InvalidTag"
	nodeWhitespace   nodeType = "Whitespace"
	nodeRoot         nodeType = "Root"
	nodeText         nodeType = "Text"
	nodeTagName      nodeType = "TagName"
	nodeBlockTag     nodeType = "BlockTag"
	nodeBlockTagLine nodeType = "BlockTagLine"
	nodeBlockTagAttr nodeType = "BlockTagAttr"
	nodeParagraph    nodeType = "Paragraph"
	nodeInlineTag    nodeType = "InlineTag"
	nodeHyperlink    nodeType = "Hyperlink"
	nodeBoldText     nodeType = "BoldText"
)

// makeNode returns a new node
// parent nodes have only a type and list of children
func (p *Parser) makeNode(n nodeType, hasVal bool) *Node {
	val := p.frame
	if !hasVal {
		val = ""
	}
	return &Node{Typ: n, Value: val, Ln: p.startLn, Col: p.startCol}
}

// createNode appends a new parent node to the node stack
func (p *Parser) createNode(n nodeType) {
	p.nodeStack = append(p.nodeStack, p.makeNode(n, false))
	p.nodeType = n
}

// currentNode returns the topmost node from the node stack
func (p *Parser) currentNode() *Node {
	if len(p.nodeStack) > 0 {
		return p.nodeStack[len(p.nodeStack)-1]
	}
	return nil
}

// setNodeType sets the type for the current node
func (p *Parser) setNodeType(n nodeType) {
	p.currentNode().Typ = n
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
func (p *Parser) addChild(n nodeType) {
	p.frame = strings.TrimSpace(p.frame)
	p.frame = strings.ReplaceAll(p.frame, "\n", " ")
	// if the last node is of the same type, add the frame content
	// to the end of the last node
	if p.frame != "" {
		if lastNode := p.lastNode(); lastNode != nil && lastNode.Typ == n {
			p.lastNode().Value += " " + p.frame
		} else {
			topNode := p.currentNode()
			topNode.Children = append(topNode.Children, p.makeNode(n, true))
		}
	}
	p.flattenFrame()
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

func (p *Parser) numChildren() int {
	return len(p.currentNode().Children)
}

func (p *Parser) parentNodeType() nodeType {
	return p.currentNode().Typ
}

func (p *Parser) determineNodeType() {
	if p.frame == "" {
		p.currentNode().Typ = nodeInlineTag
		return
	}
	var t nodeType
	switch p.frame {
	case "l":
		t = nodeHyperlink
	case "b":
		t = nodeBoldText
	default:
		t = nodeInvalidTag
		p.addError(errInvalidTagName)
	}
	p.currentNode().Typ = t
	p.nodeType = t
}
