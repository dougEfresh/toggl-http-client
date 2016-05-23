// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package gproject provides access to toggl REST API


Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/gproject"

	func main() {
		thc, err := gtoggl.NewClient("token")
		...
		pc, err := gproject.NewClient(thc)
		...
		projects,err := pc.List()
		if err == nil {
			panic(err)
		}
	}
*/
package gproject
