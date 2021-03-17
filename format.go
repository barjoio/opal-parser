package opalparser

import (
	"encoding/json"
	"fmt"
)

func (p *Parser) HTML() string {
	var html string
	for _, node := range p.tree[0].Children {
		switch node.Typ {
		case nodeTitle:
			html += "<div class='opal_Title'>\n"
			html += "\t" + htmlText(node)
			html += "</div>\n"
		case nodeToC:
		case nodeHeading:
			html += "<h" + node.Level + " class='opal_Heading'>\n"
			html += "\t" + htmlText(node)
			html += "</h" + node.Level + ">\n"
		case nodeParagraph:
			html += "<p class='opal_P'>\n"
			html += "\t" + htmlText(node)
			html += "</p>\n"
		case nodeTable:
			var hasHeader bool
			var t string
			html += "<table class='opal_Table'>\n"
			for _, attr := range node.Attrs {
				switch attr {
				case "h":
					hasHeader = true
				}
			}
			for i, row := range node.Children {
				html += "\t<tr class='opal_TableRow'>\n"
				for _, data := range row.Children {
					if i == 0 && hasHeader {
						t = "th"
					} else {
						t = "td"
					}
					html += "\t\t<" + t + " class='opal_TableData'>\n"
					html += "\t\t\t" + htmlText(data)
					html += "\t\t</" + t + ">\n"
				}
				html += "\t</tr>\n"
			}
			html += "</table>\n"
		case nodeList:
			var listType string
			if len(node.Attrs) > 0 {
				switch node.Attrs[0] {
				case "n", "number":
					html += "<ol class='opal_ListN'>\n"
					listType = "ol"
				default:
					html += "<ul class='opal_ListB'>\n"
					listType = "ul"
				}
			} else {
				html += "<ul class='opal_ListB'>\n"
				listType = "ul"
			}
			for _, listItem := range node.Children {
				html += "\t<li class='opal_ListItem'>\n"
				html += "\t\t" + htmlText(listItem)
				html += "\t</li>\n"
			}
			html += "</" + listType + ">\n"
		}
	}
	if len(html) > 0 {
		return html[:len(html)-1]
	}
	return ""
}

func bind(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func htmlText(n *Node) string {
	var html string
	for _, v := range n.Children {
		switch v.Typ {
		case nodeText:
			html += bind(" <span class='opal_Text'>%s</span>", v.Value)
		case nodeBoldText:
			html += bind(" <b class='opal_Bold'>%s</b>", v.Value)
		case nodeCode:
			html += bind(" <pre class='opal_Code'>%s</pre>", v.Value)
		case nodeHyperlink:
			html += bind(" <a class='opal_A' href='%s'>%s</a>", v.URL, v.DisplayText)
		case nodeItalicText:
			html += bind(" <i class='opal_Italic'>%s</i>", v.Value)
		case nodeUnderlineText:
			html += bind(" <u class='opal_Underline'>%s</u>", v.Value)
		}
	}
	if len(html) > 0 {
		html = html[1:]
	}
	return html + "\n"
}

func (p *Parser) JSON() string {
	b, err := json.MarshalIndent(p.tree, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
