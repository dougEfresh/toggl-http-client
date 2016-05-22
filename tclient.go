package gtoggl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DefaultAuthPassword = "api_token"
	DefaultMaxRetries   = 5
	DefaultGzipEnabled  = false
	DefaultUrl          = "https://www.toggl.com/api/v8"
	DefaultVersion      = "v8"
)

// Client is an Toggl REST client. Create one by calling NewClient.
type TogglClient struct {
	client        *http.Client // net/http Client to use for requests
	version       string       // v8
	url           string       // set of URLs passed initially to the client
	errorLog      Logger       // error log for critical messages
	infoLog       Logger       // information log for e.g. response times
	traceLog      Logger       // trace log for debugging
	password      string       // password for HTTP Basic Auth
	maxRetries    uint
	sessionCookie string //24 hour session cookie
	gzipEnabled   bool   // gzip compression enabled or disabled (default)
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*TogglClient) error

// An error is also returned when some configuration option is invalid
//    tc,err := gtoggl.NewClient("token")
func NewClient(key string, options ...ClientOptionFunc) (*TogglClient, error) {
	// Set up the client
	c := &TogglClient{
		client:      http.DefaultClient,
		maxRetries:  DefaultMaxRetries,
		url:         DefaultUrl,
		version:     DefaultVersion,
		gzipEnabled: DefaultGzipEnabled,
		password:    DefaultAuthPassword,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if len(key) < 1 {
		return nil, errors.New("Token required")
	}

	_, err := c.authenticate(key)

	if err != nil {
		return nil, err
	}

	return c, nil
}

// SetHttpClient can be used to specify the http.Client to use when making
// HTTP requests to Toggl
func SetHttpClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *TogglClient) error {
		if httpClient != nil {
			c.client = httpClient
		} else {
			c.client = http.DefaultClient
		}
		return nil
	}
}

// SetURL defines the base URL. See DefaultUrl
func SetURL(url string) ClientOptionFunc {
	return func(c *TogglClient) error {
		switch len(url) {
		case 0:
			c.url = DefaultUrl
		default:
			c.url = url
		}
		return nil
	}
}

//Custom logger to print HTTP requests
func SetTraceLogger(l Logger) ClientOptionFunc {
	return func(c *TogglClient) error {
		c.traceLog = l
		return nil
	}
}

//Custom logger to handle error messages
func SetErrorLogger(l Logger) ClientOptionFunc {
	return func(c *TogglClient) error {
		c.errorLog = l
		return nil
	}
}

//Custom logger to handle info messages
func SetInfoLogger(l Logger) ClientOptionFunc {
	return func(c *TogglClient) error {
		c.infoLog = l
		return nil
	}
}

// Returns the session cookie
func (c *TogglClient) String() string {
	return fmt.Sprintf("{sessionCookie=%s}", c.sessionCookie)
}

func (c *TogglClient) authenticate(key string) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.url, "sessions"), nil)
	if err != nil {
		return nil, err
	}
	c.dumpRequest(req)
	req.SetBasicAuth(key, "api_token")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.dumpResponse(resp)
	cookies := resp.Cookies()
	for _, value := range cookies {
		if value.Name == "toggl_api_session_new" {
			c.sessionCookie = value.Value
		}
	}

	defer resp.Body.Close()
	if resp.Body != nil {
		return ioutil.ReadAll(resp.Body)
	}
	return nil, nil
}

func request(c *TogglClient, method, endpoint string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	c.dumpRequest(req)
	cookie := &http.Cookie{}
	cookie.Name = "toggl_api_session_new"
	cookie.Value = c.sessionCookie
	req.AddCookie(cookie)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.dumpResponse(resp)

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Utility to POST requests
func (c *TogglClient) PostRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "POST", endpoint, body)
}

// Utility to DELETE requests
func (c *TogglClient) DeleteRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "DELETE", endpoint, body)
}

// Utility to PUT requests
func (c *TogglClient) PutRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "PUT", endpoint, body)
}

// Utility to GET requests
func (c *TogglClient) GetRequest(endpoint string) ([]byte, error) {
	return request(c, "GET", endpoint, nil)
}
