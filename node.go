package opalparser2

// Node is a grammatically defined element in the Opal language
// these are used to construct the abstract syntax tree
type Node struct {
	Type     nodeType `json:"type,omitempty"`
	Value    string   `json:"value,omitempty"`
	Ln       int      `json:"line,omitempty"`
	Col      int      `json:"column,omitempty"`
	Children []*Node  `json:"children,omitempty"`
}

type nodeType string

// list of nodes that can be parsed
const (
	nodeEOF           nodeType = "EOF"
	nodeWhitespace    nodeType = "Whitespace"
	nodeParagraph     nodeType = "Paragraph"
	nodeParagraphText nodeType = "ParagraphText"
	nodeInlineTag     nodeType = "InlineTag"
	nodeInlineTagName nodeType = "InlineTagName"
	nodeInlineTagText nodeType = "InlineTagText"
)

// createNode returns a new node
// parent nodes have only a type and list of children
func (p *Parser) createNode(n nodeType, isParentNode bool) *Node {
	if !isParentNode {
		return &Node{n, p.frame, p.startLn, p.startCol, []*Node{}}
	}
	return &Node{n, "", 0, 0, []*Node{}}
}

// createParentNode appends a new parent node to the node stack
func (p *Parser) createParentNode(n nodeType) {
	p.nodeStack = append(p.nodeStack, p.createNode(n, true))
}

// topNode returns the topmost node from the node stack
func (p *Parser) topNode() *Node {
	if len(p.nodeStack) > 0 {
		return p.nodeStack[len(p.nodeStack)-1]
	}
	return nil
}

// addChild appends a new child node to the topmost node from the node stack
func (p *Parser) addChild(n nodeType) {
	p.trimSpace()
	topNode := p.topNode()
	if p.frame != "" {
		topNode.Children = append(topNode.Children, p.createNode(n, false))
	}
	p.flattenFrame()
}

// popNode removes the topmost node from the node stack
func (p *Parser) popNode() {
	p.nodeStack = p.nodeStack[:len(p.nodeStack)-1]
}

// addToParent adds the topmost node from the node stack to the abstract syntax tree and removes it from the node stack
func (p *Parser) addToParent() {
	topNode := p.topNode()
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
