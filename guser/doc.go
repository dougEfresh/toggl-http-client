// Copyright 2016 Douglas Chimento.  All rights reserved.

/*
Package guser provides access to toggl REST API


Example:
	import "github/dougEfresh/gtoggl"
	import "github/dougEfresh/gtoggl/guser"

	func main() {
		thc, err := gtoggl.NewClient("token")
		...
		uc, err := guser.NewClient(thc)
		...
		users,err := uc.List()
		if err == nil {
			panic(err)
		}
	}
*/
package guser
