package gtoggl

import (
	"net/http"
	"testing"
)

func mockMe() mockFunc {
	return func(req *http.Request) string {
		return ""
	}
}

func TestClientDefaults(t *testing.T) {
	client, err := NewClient("")
	if err == nil {
		t.Fatal("Error should have been thrown. No Token given")
	}
	httpClient := &http.Client{Transport: newMockTransport(mockMe())}
	client, err = NewClient("abc1234567890def", SetURL("https://blah"), SetErrorLogger(testLogger),SetTraceLogger(testLogger),SetInfoLogger(testLogger), SetHttpClient(httpClient))
	if err != nil {
		panic(err)
	}
	if client.url != "https://blah" {
		t.Error("Url not blah; "+ client.url)
	}
	if client.traceLog != testLogger || client.errorLog != testLogger || client.infoLog != testLogger {
		t.Error("Loggers not set ")
	}
	if len(client.sessionCookie) < 1 {
		t.Errorf("Token not defined %d", len(client.sessionCookie))
	}
}
