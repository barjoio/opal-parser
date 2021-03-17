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
Bold italic text | `c Lorem \`bi ipsum\`.`
Hyperlinks | `c Lorem \`l ipsum example.com\`. or \`l _ example.com\`.`

.2: Block elements

.table/h
Name | Example | Attributes
Table of contents | `c .Toc` | None
Headings | `c .1, .2, .3, .4, .5` | None
Table | `c .table` | `c h, f`
List | `c .list` | `c b, n`
