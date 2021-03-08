package main

import (
	opalparser2 "github.com/barjoco/opalparser"
)

func main() {
	p := opalparser2.New()

	p.ParseFile("example/test.opal")
}
