package gtoggl

import (
	"encoding/json"
	"fmt"
)

type WorkspaceLister interface {
	List(tc *TogglClient) ([]Workspace, error)
}
type WorkspaceGetter interface {
	Get(tc *TogglClient, id uint64, wtype string) (Workspace, error)
}
type WorkspaceUpdater interface {
	Update(tc *TogglClient, ws Workspace) (Workspace, error)
}

func (wl *WorkspaceTransport) Get(tc *TogglClient, id uint64, wtype string) (Workspace, error) {
	body, err := tc.GetRequest(fmt.Sprintf("%s/%d", workspaceUrl, id))
	if err != nil {
		return Workspace{}, err
	}

	var aux workspace_response
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

func (wl *WorkspaceTransport) List(tc *TogglClient) ([]Workspace, error) {
	body, err := tc.GetRequest(workspaceUrl)
	var workspaces []Workspace
	if err != nil {
		return workspaces, err
	}
	err = json.Unmarshal(body, &workspaces)
	return workspaces, err
}

func (wl *WorkspaceTransport) Update(tc *TogglClient, ws Workspace) (Workspace, error) {
	put := workspace_update_request{Workspace: ws}
	body, err := json.Marshal(put)
	if err != nil {
		return Workspace{}, err
	}
	response, err := tc.PutRequest(fmt.Sprintf("%s/%d", workspaceUrl, ws.Id), body)
	if err != nil {
		return Workspace{}, err
	}
	var aux workspace_response
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

type WorkspaceClient struct {
	tc              *TogglClient
	listTransport   WorkspaceLister
	getTransport    WorkspaceGetter
	updateTransport WorkspaceUpdater
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type WorkspaceClientOptionFunc func(*WorkspaceClient) error
type WorkspaceTransport struct {
}

var defaultTransport = &WorkspaceTransport{}

func NewWorkspaceClient(tc *TogglClient, options ...WorkspaceClientOptionFunc) (*WorkspaceClient, error) {
	ws := &WorkspaceClient{
		tc: tc,
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
	return ws, nil
}

func SetTogglClient(tc *TogglClient) WorkspaceClientOptionFunc {
	return func(ws *WorkspaceClient) error {
		ws.tc = tc
		return nil
	}
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

type workspace_response struct {
	Data Workspace `json:"data"`
}

type workspace_update_request struct {
	Workspace Workspace `json:"workspace"`
}

var workspaceUrl = DefaultUrl + "/workspaces"

func (wc *WorkspaceClient) Get(id uint64) (Workspace, error) {
	return wc.getTransport.Get(wc.tc, id, "")
}

func (wc *WorkspaceClient) Update(ws Workspace) (Workspace, error) {
	return wc.updateTransport.Update(wc.tc, ws)
}

func (wc *WorkspaceClient) List() ([]Workspace, error) {
	return wc.listTransport.List(wc.tc)
}

func (ws *WorkspaceClient) String() string {
	return fmt.Sprintf("workspace:{togglClient: %s}", ws.tc)
}
