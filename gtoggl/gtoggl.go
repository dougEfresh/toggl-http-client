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
	if *command == "client" {
		client(tc, flag.Args())
	}
}

func handleError(error error) {
	if error != nil {
		fmt.Fprintln(os.Stderr, error)
		os.Exit(-1)
	}
}

func client(tc *gtoggl.TogglHttpClient, args []string) {
	c, err := gtoggl.NewTogglClient(tc)
	var client gtoggl.Client
	handleError(err)
	if len(args) == 0 || args[0] == "list" {
		clients, err := c.List()
		handleError(err)
		fmt.Printf("%+v\n", clients)
	}

	if args[0] == "create" && len(args) > 1 {
		err = json.Unmarshal([]byte(args[1]), &client)
		handleError(err)
		client, err = c.Create(&client)
		handleError(err)
		fmt.Printf("%+v\n", client)
	}

	if args[0] == "update" && len(args) > 1 {
		err = json.Unmarshal([]byte(args[1]), &client)
		handleError(err)
		_, err = c.Get(client.Id)
		handleError(err)
		client, err = c.Update(&client)
		handleError(err)
		fmt.Printf("%+v\n", client)
	}

	if args[0] == "get" && len(args) > 1 {
		i, err := strconv.ParseUint(args[1], 0, 64)
		handleError(err)
		client, err = c.Get(i)
		handleError(err)
		fmt.Printf("%+v\n", client)
	}

	if args[0] == "delete" && len(args) > 1 {
		i, err := strconv.ParseUint(args[1], 0, 64)
		handleError(err)
		err = c.Delete(i)
		handleError(err)
		fmt.Printf("%+v  deleted\n", i)
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
		var uWs gtoggl.Workspace
		handleError(err)
		err = json.Unmarshal([]byte(args[1]), &uWs)
		handleError(err)
		_, err := wsc.Get(uWs.Id)
		handleError(err)
		uWs, err = wsc.Update(uWs)
		handleError(err)
		fmt.Printf("%+v\n", uWs)
		return
	}
}
