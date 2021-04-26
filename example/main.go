package main

import (
	"github.com/barjoco/opalparser"
)

func main() {
	p := opalparser.New()

	p.ParseFile("example/test.opal")

	pdf := p.PDF()
	pdf.SaveAs("example/test.pdf")
}
