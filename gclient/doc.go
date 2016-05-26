// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package gclient provides access to toggl REST API.

Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/gclient"

	func main() {
		thc, err := gtoggl.NewClient("token")
		...
		tc, err := gclient.NewClient(thc)
		...
		clients,err := tc.List()
		if err == nil {
			panic(err)
		}
	}
*/
package gclient
