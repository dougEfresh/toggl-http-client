package gtoggl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type TestLogger struct {
}

func (l *TestLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

var l = &TestLogger{}

type mockTransport struct{}

func newMockTransport() http.RoundTripper {
	return &mockTransport{}
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create mocked http.Response
	response := &http.Response{Header: make(http.Header), Request: req, StatusCode: http.StatusOK}
	response.Header.Set("Content-Type", "application/json")
	responseBody := GetResponse(req)
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

func GetResponse(req *http.Request) string {
	r := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	if r == "GET /api/v8/workspaces" {
		b, err := ioutil.ReadFile("mock/workspaces.json")
		if err != nil {
			panic(err)
		}
		return string(b)
	}

	if r == "GET /api/v8/workspaces/1" {
		b, err := ioutil.ReadFile("mock/workspace.json")
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	if r == "PUT /api/v8/workspaces/1" {
		b, err := ioutil.ReadFile("mock/workspace_update.json")
		if err != nil {
			panic(err)
		}
		return string(b)
	}

	panic(errors.New("Cannot mock an unknown request"))
}

func workspaceClient() *WorkspaceClient {
	httpClient := &http.Client{Transport: newMockTransport()}
	client, err := NewClient("abc1234567890def", SetTraceLogger(l), SetHttpClient(httpClient))
	if err != nil {
		panic(err)
	}
	ws, err := NewWorkspaceClient(SetTogglClient(client))
	if err != nil {
		panic(err)
	}
	return ws
}

func TestWorkspaceList(t *testing.T) {
	workspaceClient := workspaceClient()
	workspaces, err := workspaceClient.List()
	if err != nil {
		t.Error(err)
	}
	if len(workspaces) != 2 {
		t.Error("Workspace is not 2")
	}
	if workspaces[0].Id != 1 {
		t.Error("Workspace Id is not 1")
	}
	if workspaces[0].Name != "Id 1" {
		t.Error("Workspace name not Id ")
	}
	if !workspaces[0].Premium {
		t.Error("Workspace is not premium ")
	}

	if workspaces[1].Id != 2 {
		t.Error("Workspace Id is not 2")
	}
	if workspaces[1].Name != "Id 2" {
		t.Error("Workspace name")
	}
	if workspaces[1].Premium {
		t.Error("Workspace is not premium ")
	}

}

func TestWorkspaceGet(t *testing.T) {
	workspaceClient := workspaceClient()

	workspace, err := workspaceClient.Get(1)
	if err != nil {
		t.Error(err)
	}
	if workspace.Id != 1 {
		t.Error("Workspace id != 1")
	}
}

func TestWorkspaceUpdate(t *testing.T) {
	workspaceClient := workspaceClient()
	workspace, err := workspaceClient.Get(1)

	if err != nil {
		t.Error(err)
	}
	workspace.Name = "new name"
	workspace, err = workspaceClient.Update(workspace)

	if err != nil {
		t.Error(err)
	}

	if workspace.Name != "new name" {
		t.Error("Wrong name")
	}

}
