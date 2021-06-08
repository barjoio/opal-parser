package opalparser

import (
	"reflect"
	"testing"
)

type ParsingTest struct {
	rawMarkup     string
	expectedNodes []nodeType
}

var parsingTests = []ParsingTest{
	{"", []nodeType{nodeRoot}},
	{"Foo bar baz", []nodeType{nodeRoot, nodeParagraph, nodeText}},
	{"Foo; bar; baz", []nodeType{nodeRoot, nodeParagraph, nodeText, nodeParagraph, nodeText, nodeParagraph, nodeText}},
	{"Foo\\; bar\\; baz", []nodeType{nodeRoot, nodeParagraph, nodeText}},
	{"Foo `b bar` baz", []nodeType{nodeRoot, nodeParagraph, nodeText, nodeInlineTag, nodeText}},
	{"Foo \\`b bar\\` baz", []nodeType{nodeRoot, nodeParagraph, nodeText}},
	{".1: Foo bar baz", []nodeType{nodeRoot, nodeBlockTag, nodeText}},
	{"\\.1: Foo bar baz", []nodeType{nodeRoot, nodeParagraph, nodeText}},
	{".1: Foo `b bar` baz", []nodeType{nodeRoot, nodeBlockTag, nodeText, nodeInlineTag, nodeText}},
	{"Foo `l bar example.com` baz", []nodeType{nodeRoot, nodeParagraph, nodeText, nodeInlineTag, nodeText}},
	{"Foo `l _ example.com` baz", []nodeType{nodeRoot, nodeParagraph, nodeText, nodeInlineTag, nodeText}},
	{".list\n" +
		"- foo\n" +
		"- bar\n" +
		"- baz", []nodeType{nodeRoot, nodeBlockTag, nodeListItem, nodeText, nodeListItem, nodeText, nodeListItem, nodeText, nodeListItem, nodeText, nodeListItem}},
	{".list\n" +
		"- foo\n" +
		"- `b bar`\n" +
		"- baz", []nodeType{nodeRoot, nodeBlockTag, nodeListItem, nodeText, nodeListItem, nodeText, nodeListItem, nodeText, nodeInlineTag, nodeText, nodeListItem, nodeText, nodeListItem}},
	{".table\n" +
		"abc | def | ghi\n" +
		"jkl | mno | pqr\n" +
		"stu | vwx | yz", []nodeType{nodeRoot, nodeBlockTag, nodeTableRow, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableRow, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableRow, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableData, nodeText, nodeTableData}},
}

func TestParse(t *testing.T) {
	for _, test := range parsingTests {
		p := New()
		p.Parse(test.rawMarkup)

		nodes := p.nodesParsedLinear

		if !reflect.DeepEqual(nodes, test.expectedNodes) {
			t.Fatalf("Expected nodes: %v, got nodes: %v", test.expectedNodes, nodes)
		}
	}
}
