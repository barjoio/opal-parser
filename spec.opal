.Title: Opal Language Specification

.ToC

.1: Paragraphs

`c ...`
`c Lorem ipsum`
`c ...`

.1: Inline elements

.table/h
Name | Example
Bold text | `c Lorem \`b ipsum\`.`
Underline text | `c Lorem \`u ipsum\`.`
Italic text | `c Lorem \`i ipsum\`.`
Hyperlinks | `c Lorem \`l ipsum example.com\`` or `c Lorem \`l _ example.com\`.`

.1: Block elements

.2: Attributes

Apply attributes to an element using the ".elementName/attribute1/attribute2" notation.

For example:

`c .list/b`
`c - item 1`
`c - item 2`
`c - item 3`

Where the "b" attribute, for "bullet", can be switched to "n" to render a numbered list.

Block elements along with their associated attributes are shown in the table below.

.table/h
Name | Example | Attributes
Headings | `c .1: Lorem ipsum` also `c .2, .3, .4, .5, .6` | None
Table | `c .table` `c abc | def | ghi` `c jkl | mno | pqr` | `c h, f`
List | `c .list` | `c b, n`

.1: Misc

.list/b
- Terminators are expressed using a semicolon (\;), or double newline (a collection of whitespace which includes 2 or more newline characters).