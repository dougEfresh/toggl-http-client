package gtoggl

import (
	"encoding/json"
	"fmt"
)

type WorkspaceLister interface {
	List(wsc *WorkspaceClient) ([]Workspace, error)
}
type WorkspaceGetter interface {
	Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error)
}
type WorkspaceUpdater interface {
	Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error)
}

func (wl *WorkspaceTransport) Get(wsc *WorkspaceClient, id uint64, wtype string) (Workspace, error) {
	body, err := wsc.tc.GetRequest(fmt.Sprintf("%s/%d",wsc.workspaceEndpoint, id))
	if err != nil {
		return Workspace{}, err
	}

	var aux workspace_response
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

func (wl *WorkspaceTransport) List(wsc *WorkspaceClient) ([]Workspace, error) {
	body, err := wsc.tc.GetRequest(wsc.workspaceEndpoint)
	var workspaces []Workspace
	if err != nil {
		return workspaces, err
	}
	err = json.Unmarshal(body, &workspaces)
	return workspaces, err
}

func (wl *WorkspaceTransport) Update(wsc *WorkspaceClient, ws Workspace) (Workspace, error) {
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

type WorkspaceClient struct {
	tc              *TogglClient
	workspaceEndpoint string
	listTransport   WorkspaceLister
	getTransport    WorkspaceGetter
	updateTransport WorkspaceUpdater
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewWorkspaceClient.
type WorkspaceClientOptionFunc func(*WorkspaceClient) error
type WorkspaceTransport struct {
}

var defaultTransport = &WorkspaceTransport{}

func NewWorkspaceClient(tc *TogglClient, options ...WorkspaceClientOptionFunc) (*WorkspaceClient, error) {
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
	ws.workspaceEndpoint = DefaultUrl + "/workspaces"
	return ws, nil
}

func SetGetTransport(g WorkspaceGetter) WorkspaceClientOptionFunc {
	return func(ws *WorkspaceClient) error {
		ws.getTransport = g
		return nil
	}
}

type Workspace struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}
type Workspaces []Workspace

func (ws *Workspace) String() string {
	return fmt.Sprintf("{id:%s ,name:%s, premium: %b",ws.Id,ws.Name,ws.Premium)
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

func (wc *WorkspaceClient) Get(id uint64) (Workspace, error) {
	return wc.getTransport.Get(wc, id, "")
}

func (wc *WorkspaceClient) Update(ws Workspace) (Workspace, error) {
	return wc.updateTransport.Update(wc, ws)
}

func (wc *WorkspaceClient) List() ([]Workspace, error) {
	return wc.listTransport.List(wc)
}

func (ws *WorkspaceClient) String() string {
	return fmt.Sprintf("workspace:{togglClient: %s}", ws.tc)
}
