package gtoggl

import (
	"encoding/json"
	"fmt"
)

type Workspace struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}

type Workspaces []Workspace

//Return a Workspace Cilent. An error is also returned when some configuration option is invalid
//    tc,err := gtoggl.NewClient("token")
//    wsc,err := gtoggl.NewWorkspaceClient(tc)
func NewWorkspaceClient(tc *TogglHttpClient, options ...WorkspaceClientOptionFunc) (*WorkspaceClient, error) {
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
	ws.workspaceEndpoint = tc.Url + "/workspaces"
	return ws, nil
}

type WorkspaceClient struct {
	tc                *TogglHttpClient
	workspaceEndpoint string
	listTransport     WorkspaceLister
	getTransport      WorkspaceGetter
	updateTransport   WorkspaceUpdater
}

func (wc *WorkspaceClient) Get(id uint64) (Workspace, error) {
	return wc.getTransport.Get(wc, id, "")
}

func (wc *WorkspaceClient) Update(ws Workspace) (Workspace, error) {
	return wc.updateTransport.Update(wc, ws)
}

func (wc *WorkspaceClient) List() (Workspaces, error) {
	return wc.listTransport.List(wc)
}

func (ws *WorkspaceClient) String() string {
	return fmt.Sprintf("workspace:{togglHttpClient: %s}", ws.tc)
}

type WorkspaceLister interface {
	List(wsc *WorkspaceClient) (Workspaces, error)
}
type WorkspaceGetter interface {
	Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error)
}
type WorkspaceUpdater interface {
	Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error)
}

func (wl *workspaceTransport) Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error) {
	body, err := wsc.tc.GetRequest(fmt.Sprintf("%s/%d", wsc.workspaceEndpoint, id))
	if err != nil {
		return Workspace{}, err
	}

	var aux workspace_response
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

func (wl *workspaceTransport) List(wsc *WorkspaceClient) (Workspaces, error) {
	body, err := wsc.tc.GetRequest(wsc.workspaceEndpoint)
	var workspaces []Workspace
	if err != nil {
		return workspaces, err
	}
	err = json.Unmarshal(body, &workspaces)
	return workspaces, err
}

func (wl *workspaceTransport) Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error) {
	put := workspace_update_request{Workspace: ws}
	body, err := json.Marshal(put)
	if err != nil {
		return Workspace{}, err
	}
	response, err := wsc.tc.PutRequest(fmt.Sprintf("%s/%d", wsc.workspaceEndpoint, ws.Id), body)
	if err != nil {
		return Workspace{}, err
	}
	var aux workspace_response
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewWorkspaceClient.
type WorkspaceClientOptionFunc func(*WorkspaceClient) error

type workspaceTransport struct {
}

var defaultTransport = &workspaceTransport{}

func (ws *Workspace) String() string {
	return fmt.Sprintf("{id:%s ,name:%s, premium: %b", ws.Id, ws.Name, ws.Premium)
}

func (ws Workspaces) String() string {
	var st string = "["
	for _, value := range ws {
		st += fmt.Sprint(value) + ","
	}
	return st + "]"
}

type workspace_response struct {
	Data Workspace `json:"data"`
}

type workspace_update_request struct {
	Workspace Workspace `json:"workspace"`
}
