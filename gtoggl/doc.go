// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package gtoggl access to toggl REST API

Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/gclient"

	func main() {
		thc, err := gtoggl.NewClient("token")
		thc, err = gtoggl.NewClient("token",gtoggl.SetURL("https://www.toggl.com/api/v8/")
		...
		...
		clients,err := tc.List()
		if err == nil {
			panic(err)
		}
	}
*/
package main
