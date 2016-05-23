/// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package gworkspace provides access to toggl REST API

Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/gworkspace"

	func main() {
		thc, err := gtoggl.NewClient("token")
		...
		wsc, err := gworkspace.NewClient(thc)
		...
		workspaces,err := wsc.List()
		if err == nil {
			panic(err)
		}
	}
*/
package gworkspace
