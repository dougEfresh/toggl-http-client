package gtoggl

import "testing"

func TestClientDefaults(t *testing.T) {
	client, err := NewClient("")
	if err == nil {
		t.Fatal("Error should have been thrown. No Token given")
	}
	client, err = NewClient("test-token")

	if client == nil {
		t.Fatal(err)
	}

	if len(client.token) < 1 {
		t.Error("Token not defined %d", len(client.token))
	}
}
