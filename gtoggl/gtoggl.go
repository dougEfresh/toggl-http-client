package main

import (
	"flag"
	"fmt"
	"github.com/dougEfresh/gtoggl"
	"os"
)

func main() {
	flag.String("t", "", "Toggl API token: https://www.toggl.com/app/profile")
	flag.Parse()

	token := flag.Lookup("t")

	tc, err := gtoggl.NewClient(token.Value.String())
	if err != nil {
		fmt.Fprint(os.Stderr,"A token is required\n")
		flag.Usage()
		os.Exit(-1)
	}
	wsc, err := gtoggl.NewWorkspaceClient(tc)
	fmt.Print(wsc.String())
}
