package gworkspace

import (
	"encoding/json"
	"fmt"
	"github.com/dougEfresh/gtoggl"
)

type Workspace struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}

type Workspaces []Workspace

const Endpoint = "/workspaces"

//Return a Workspace Cilent. An error is also returned when some configuration option is invalid
//    tc,err := gtoggl.NewClient("token")
//    wsc,err := gtoggl.NewWorkspaceClient(tc)
func NewClient(tc *gtoggl.TogglHttpClient, options ...WorkspaceClientOptionFunc) (*WorkspaceClient, error) {
	ws := &WorkspaceClient{
		tc:              tc,
		listTransport:   defaultTransport,
		getTransport:    defaultTransport,
		updateTransport: defaultTransport,
	}
	// Run the options on it

	for _, option := range options {
		if err := option(ws); err != nil {
			return nil, err
		}
	}
	ws.endpoint = tc.Url + Endpoint
	return ws, nil
}

type WorkspaceClient struct {
	tc              *gtoggl.TogglHttpClient
	endpoint        string
	listTransport   WorkspaceLister
	getTransport    WorkspaceGetter
	updateTransport WorkspaceUpdater
}

//https://github.com/toggl/toggl_api_docs/blob/master/chapters/workspaces.md
func (wc *WorkspaceClient) Get(id uint64) (Workspace, error) {
	return wc.getTransport.Get(wc, id, "")
}

//https://github.com/toggl/toggl_api_docs/blob/master/chapters/workspaces.md
func (wc *WorkspaceClient) Update(ws Workspace) (Workspace, error) {
	return wc.updateTransport.Update(wc, ws)
}

//https://github.com/toggl/toggl_api_docs/blob/master/chapters/workspaces.md
func (wc *WorkspaceClient) List() (Workspaces, error) {
	return wc.listTransport.List(wc)
}

//Ability to override the List method. Not yet implemented
type WorkspaceLister interface {
	List(wsc *WorkspaceClient) (Workspaces, error)
}

//Ability to override the Get method. Not yet implemented
type WorkspaceGetter interface {
	Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error)
}

//Ability to override the Update method. Not yet implemented
type WorkspaceUpdater interface {
	Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error)
}

//GET https://www.toggl.com/api/v8/workspaces/123213
func (wl *workspaceTransport) Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error) {
	body, err := wsc.tc.GetRequest(fmt.Sprintf("%s/%d", wsc.endpoint, id))
	if err != nil {
		return Workspace{}, err
	}

	var aux workspace_response
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

//GET https://www.toggl.com/api/v8/workspaces
func (wl *workspaceTransport) List(wsc *WorkspaceClient) (Workspaces, error) {
	body, err := wsc.tc.GetRequest(wsc.endpoint)
	var workspaces []Workspace
	if err != nil {
		return workspaces, err
	}
	err = json.Unmarshal(body, &workspaces)
	return workspaces, err
}

//PUT https://www.toggl.com/api/v8/workspaces
func (wl *workspaceTransport) Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error) {
	put := workspace_update_request{Workspace: ws}
	body, err := json.Marshal(put)
	if err != nil {
		return Workspace{}, err
	}
	response, err := wsc.tc.PutRequest(fmt.Sprintf("%s/%d", wsc.endpoint, ws.Id), body)
	if err != nil {
		return Workspace{}, err
	}
	var aux workspace_response
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

//Configures a Client.
/*
    func SetURL(url string) WorkspaceClientOptionFunc {
	return func(c *WorkspaceClient) error {
	    c.Url = url
	}
    }
*/
type WorkspaceClientOptionFunc func(*WorkspaceClient) error

type workspaceTransport struct {
}

var defaultTransport = &workspaceTransport{}

type workspace_response struct {
	Data Workspace `json:"data"`
}

type workspace_update_request struct {
	Workspace Workspace `json:"workspace"`
}
