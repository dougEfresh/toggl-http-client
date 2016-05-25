package gclient

import (
	"github.com/dougEfresh/gtoggl/test"
	"os"
	"testing"
)

var _, debugMode = os.LookupEnv("GTOGGL_TEST_DEBUG")

func togglClient() *TogglClient {
	tu := &gtest.TestUtil{Debug: debugMode}
	client := tu.MockClient()
	ws, err := NewClient(client)
	if err != nil {
		panic(err)
	}
	return ws
}

func TestClientCreate(t *testing.T) {
	tClient := togglClient()
	c := &Client{Name: "Very Big Company", WId: 777}
	nc, err := tClient.Create(c)
	if err != nil {
		t.Fatal(err)
	}

	if nc.Name != "Very Big Company" {
		t.Fatal("!= Very Big Company")
	}

	if nc.Id != 1239455 {
		t.Fatal("!= 1239455")
	}

	if nc.WId != 777 {
		t.Fatal("!= 777")
	}
}

func TestClientUpdate(t *testing.T) {
	tClient := togglClient()
	c := &Client{Id: 1, Name: "new name", WId: 777}
	nc, err := tClient.Update(c)
	if err != nil {
		t.Fatal(err)
	}

	if nc.Name != "new name" {
		t.Fatal("!= new name")
	}
}

func TestClientDelete(t *testing.T) {
	tClient := togglClient()
	c := &Client{Id: 1, Name: "new name", WId: 777}
	err := tClient.Delete(c.Id)
	if err != nil {
		t.Fatal(err)
	}
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

func BenchmarkClientTransport_Get(b *testing.B) {
	b.ReportAllocs()
	tClient := togglClient()
	for i := 0; i < b.N; i++ {
		tClient.Get(1)
	}
}

func BenchmarkClientTransport_List(b *testing.B) {
	b.ReportAllocs()
	tClient := togglClient()
	for i := 0; i < b.N; i++ {
		tClient.List()
	}
}
