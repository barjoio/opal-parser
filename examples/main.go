package main

import (
	"encoding/json"
	"fmt"

	"github.com/barjoco/opalparser"
)

func main() {
	op := opalparser.New()

	j, err := json.MarshalIndent(op.ParseFile("examples/file.opal"), "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
