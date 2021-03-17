package main

import (
	"fmt"

	opalparser2 "github.com/barjoco/opalparser"
)

func main() {
	p := opalparser2.New()

	p.ParseFile("example/test.opal")

	h := p.HTML()

	fmt.Println(h)
}
