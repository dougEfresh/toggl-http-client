package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dougEfresh/gtoggl"
	"os"
	"strconv"
)

type Debugger struct {
	debug bool
}

func (l *Debugger) Printf(format string, v ...interface{}) {
	if l.debug {
		fmt.Printf(format, v)
	}
}

func main() {
	var debug = flag.Bool("d", false, "Debug")
	var token = flag.String("t", "", "Toggl API token: https://www.toggl.com/app/profile")
	var command = flag.String("c", "workspace", "Sub command: workspace,client,project...etc ")
	flag.Parse()
	tc, err := gtoggl.NewClient(*token, gtoggl.SetTraceLogger(&Debugger{debug: *debug}))
	if err != nil {
		fmt.Fprint(os.Stderr, "A token is required\n")
		flag.Usage()
		os.Exit(-1)
	}
	if *command == "workspace" {
		workspace(tc, flag.Args())
	}
}

func handleError(error error) {
	if error != nil {
		fmt.Fprint(os.Stderr, error)
		os.Exit(-1)
	}
}

func workspace(tc *gtoggl.TogglHttpClient, args []string) {
	wsc, err := gtoggl.NewWorkspaceClient(tc)
	handleError(err)
	if len(args) == 0 || args[0] == "list" {
		w, err := wsc.List()
		handleError(err)
		fmt.Printf("%+v\n", w)
		return
	}

	if args[0] == "get" && len(args) > 1 {
		i, err := strconv.ParseUint(args[1], 0, 64)
		handleError(err)
		w, err := wsc.Get(i)
		fmt.Printf("%+v\n", w)
		return
	}
	if args[0] == "update" && len(args) > 1 {
		i, err := strconv.ParseUint(args[1], 0, 64)
		handleError(err)
		ws, err := wsc.Get(i)
		var uWs = gtoggl.Workspace{Id: ws.Id, Name: ws.Name, Premium: ws.Premium}
		handleError(err)
		err = json.Unmarshal([]byte(args[2]), &uWs)
		handleError(err)
		uWs, err = wsc.Update(uWs)
		handleError(err)
		fmt.Printf("%+v\n", uWs)
		return
	}
}
