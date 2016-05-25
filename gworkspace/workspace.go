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
		tc: tc,
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
	tc       *gtoggl.TogglHttpClient
	endpoint string
}

//GET https://www.toggl.com/api/v8/workspaces/123213
func (wc *WorkspaceClient) Get(id uint64) (*Workspace, error) {
	return workspaceResponse(wc.tc.GetRequest(fmt.Sprintf("%s/%d", wc.endpoint, id)))
}

//PUT https://www.toggl.com/api/v8/workspaces
func (wc *WorkspaceClient) Update(ws *Workspace) (*Workspace, error) {
	put := workspace_update_request{Workspace: ws}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return workspaceResponse(wc.tc.PutRequest(fmt.Sprintf("%s/%d", wc.endpoint, ws.Id), body))
}

//GET https://www.toggl.com/api/v8/workspaces
func (wc *WorkspaceClient) List() (Workspaces, error) {
	body, err := wc.tc.GetRequest(wc.endpoint)
	var workspaces Workspaces
	if err != nil {
		return workspaces, err
	}
	err = json.Unmarshal(*body, &workspaces)
	return workspaces, err
}

func workspaceResponse(response *json.RawMessage, error error) (*Workspace, error) {
	if error != nil {
		return nil, error
	}
	var tResp gtoggl.TogglResponse
	var ws Workspace
	err := json.Unmarshal(*response, &tResp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*tResp.Data, &ws)
	if err != nil {
		return nil, err
	}
	return &ws, err
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

type workspace_update_request struct {
	Workspace *Workspace `json:"workspace"`
}
