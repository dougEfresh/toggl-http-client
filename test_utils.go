package gtoggl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TestLogger struct {
}

func (l *TestLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

var testLogger = &TestLogger{}

type mockFunc func(req *http.Request) string

type mockTransport struct {
	mock mockFunc
}

func newMockTransport(f mockFunc) http.RoundTripper {
	return &mockTransport{mock: f}
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create mocked http.Response
	response := &http.Response{Header: make(http.Header), Request: req, StatusCode: http.StatusOK}
	response.Header.Set("Content-Type", "application/json")
	response.Header.Set("Set-Cookie", "toggl_api_session_new=MTM2MzA4MJa8jA3OHxEdi1CQkFFQ180SUFBUkFCRUFBQVlQLUNBQUVHYzNSeWFXNW5EQXdBQ25ObGMzTnBiMjVmYVdRR2MzUnlhVzVuREQ0QVBIUnZaMmRzTFdGd2FTMXpaWE56YVc5dUxUSXRaalU1WmpaalpEUTVOV1ZsTVRoaE1UaGhaalpqWkRkbU5XWTJNV0psWVRnd09EWmlPVEV3WkE9PXweAkG7kI6NBG-iqvhNn1MSDhkz2Pz_UYTzdBvZjCaA==; Path=/; Expires=Wed, 13 Mar 2013 09:54:38 UTC; Max-Age=86400; HttpOnly")
	r := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	fmt.Println("DEBUG -------------- " + r)
	if strings.Contains(r,"/sessions") {
		response.Body = ioutil.NopCloser(strings.NewReader(""))
		return response, nil
	}
	responseBody := t.mock(req)
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

func mockClient(m mockFunc) *TogglClient {
	httpClient := &http.Client{Transport: newMockTransport(m)}
	optionsWithClient  := []ClientOptionFunc{SetHttpClient(httpClient)}
	client, err := NewClient("abc1234567890def", optionsWithClient...)
	if err != nil {
		panic(err)
	}
	return client
}


func mockClientOptions(m mockFunc, options []ClientOptionFunc) *TogglClient {
	httpClient := &http.Client{Transport: newMockTransport(m)}
	optionsWithClient  := options[0:len(options)+1]
	optionsWithClient[len(options)] = SetHttpClient(httpClient)
	client, err := NewClient("abc1234567890def", optionsWithClient...)
	if err != nil {
		panic(err)
	}
	return client
}
