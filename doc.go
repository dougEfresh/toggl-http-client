// Gtoggl - A go client library for toggl
/*
Package provides access to toggl REST API
https://github.com/toggl/toggl_api_docs

Example:
	import "github/dougEfresh/gtoggl"

		func main() {
			thc, err := NewClient("token")
			...
			wsc, err := NewWorkspaceClient("token")
		...
			workspaces,err := wsc.List()
			if err == nil {
			panic(err)
			}
		}
*/
package gtoggl
