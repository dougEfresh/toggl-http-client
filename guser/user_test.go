package guser

import (
	"github.com/dougEfresh/gtoggl/test"
	"os"
	"testing"
)

var _, debugMode = os.LookupEnv("GTOGGL_TEST_DEBUG")

func togglClient() *UserClient {
	tu := &gtest.TestUtil{Debug: debugMode}
	client := tu.MockClient()
	ws, err := NewClient(client)
	if err != nil {
		panic(err)
	}
	return ws
}

func TestUserCreate(t *testing.T) {
	tClient := togglClient()
	nc, err := tClient.Create("signup@blah.com", "StrongPasswrod", "UTC")
	if err != nil {
		t.Fatal(err)
	}

	if nc.Id != 3 {
		t.Fatal("!= 3")
	}

	if nc.Email != "signup@blah.com" {
		t.Fatal("!= signup@blah.com")
	}

	if nc.ApiToken != "808lolae4eab897cce9729a53642124effe" {
		t.Fatal("!= 808lolae4eab897cce9729a53642124effe")
	}
}

func TestUserUpdate(t *testing.T) {
	tClient := togglClient()
	c := &User{Id: 1, FullName: "John Swift", Email: "newemail@swift.com"}
	nc, err := tClient.Update(c)
	if err != nil {
		t.Fatal(err)
	}

	if nc.Email != "newemail@swift.com" {
		t.Fatal("!= newemail@swift.com")
	}
}

func TestUserReset(t *testing.T) {
	tClient := togglClient()
	token, err := tClient.ResetToken()
	if err != nil {
		t.Fatal(err)
	}

	if token != "123456789" {
		t.Fatal("!= 123456789")
	}
}

func TestUserGet(t *testing.T) {
	tClient := togglClient()

	client, err := tClient.Get(false)
	if err != nil {
		t.Fatal(err)
	}
	if client.Id != 1 {
		t.Error("!= 1")
	}

	if client.FullName != "John Swift" {
		t.Error("!= John Swift:  " + client.FullName)
	}

	if client.ApiToken != "1971800d4d82861d8f2c1651fea4d212" {
		t.Error("!= J1971800d4d82861d8f2c1651fea4d212:  " + client.ApiToken)
	}

	if client.Email != "johnt@swift.com" {
		t.Error("!= johnt@swift.com" + client.Email)
	}
}

func BenchmarkClientTransport_Get(b *testing.B) {
	b.ReportAllocs()
	tClient := togglClient()
	for i := 0; i < b.N; i++ {
		tClient.Get(false)
	}
}
