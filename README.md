# opalparser

A performant recursive descent parser for the Opal markup language.

## Example

```go
// New parser
p := opalparser.New()

// Path to Opal file
p.ParseFile("example/test.opal")

// Export as PDF, JSON, or HTML
pdf := p.PDF()
pdf.SaveAs("example/test.pdf")
```
