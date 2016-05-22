package gtoggl

import (
	"testing"
)

func togglClient() *TogglClient {
	client := mockClient()
	ws, err := NewTogglClient(client)
	if err != nil {
		panic(err)
	}
	return ws
}

func TestClientList(t *testing.T) {
	tClient := togglClient()
	clients, err := tClient.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(clients) != 2 {
		t.Fatal("Workspace is not 2")
	}
	if clients[0].Id != 1 {
		t.Error("Workspace Id is not 1")
	}
	if clients[0].Name != "Id 1" {
		t.Error("Workspace name not Id ")
	}

	if clients[1].Id != 2 {
		t.Error("Workspace Id is not 2")
	}
	if clients[1].Name != "Id 2" {
		t.Error("Workspace name")
	}
}

func TestClientGet(t *testing.T) {
	tClient := togglClient()

	client, err := tClient.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if client.Id != 1 {
		t.Error("!= 1")
	}

	if client.Name != "Id 1" {
		t.Error("!= Id 1:  " + client.Name)
	}
}