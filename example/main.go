package main

import (
	"fmt"

	"github.com/barjoio/opalparser"
)

func main() {
	p := opalparser.New()

	p.ParseFile("example/test.opal")

	fmt.Println(p.HTML())
}
