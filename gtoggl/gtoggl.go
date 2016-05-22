package main

import (
	"flag"
	"fmt"
	"github.com/dougEfresh/gtoggl"
)

func main() {
	flag.String("t", "", "Toggl API token: ")
	flag.Parse()
	token := flag.Lookup("t")

	c, err := gtoggl.NewClient(token.Value.String())
	if err != nil {
		panic(err)
	}

	fmt.Print(c.String())
}
