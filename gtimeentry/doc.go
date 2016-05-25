// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package gtimeentry provides access to toggl REST API


Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/gtimeentry"

	func main() {
		thc, err := gtoggl.NewClient("token")
		...
		tc, err := gtimeentry.NewClient(thc)
		...
		timeentry,err := tc.Get(1)
		if err == nil {
			panic(err)
		}
	}
*/
package gtimeentry
