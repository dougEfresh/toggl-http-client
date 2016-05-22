package gtoggl

import (
	"testing"
	"net/http"
)

func mockMe() mockFunc {
	return func (req *http.Request) string {
		return ""
	}
}

func TestClientDefaults(t *testing.T) {
	client , err := NewClient("")
	if err == nil {
		t.Fatal("Error should have been thrown. No Token given")
	}
	client = mockClient(mockMe())

	if len(client.sessionCookie) < 1 {
		t.Errorf("Token not defined %d", len(client.sessionCookie))
	}
}
