package main

import (
	"flag"
	"fmt"
	"github.com/dougEfresh/gtoggl"
	"os"
)

func main() {
	flag.Bool("d", false, "Debug")
	flag.String("t", "", "Toggl API token: https://www.toggl.com/app/profile")
	flag.String("c", "workspace", "Sub command: workspace,client,project...etc ")
	flag.Parse()
	token := flag.Lookup("t")

	tc, err := gtoggl.NewClient(token.Value.String())
	if err != nil {
		fmt.Fprint(os.Stderr, "A token is required\n")
		flag.Usage()
		os.Exit(-1)
	}
	workspace(tc)
}

func workspace(tc *gtoggl.TogglClient) {
	wsc, err := gtoggl.NewWorkspaceClient(tc)
	if err != err {

	}

	fmt.Println(wsc.List())
}
