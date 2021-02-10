package main

import (
	"fmt"
	"io/ioutil"

	"github.com/barjoco/opalparser"
)

func main() {
	b, _ := ioutil.ReadFile("examples/file.opal")
	p := opalparser.New()
	t := p.Parse(string(b))
	fmt.Println(t)
}
