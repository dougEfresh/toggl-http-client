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

// Client is an Elasticsearch client. Create one by calling NewClient.
type TogglClient struct {
	client        *http.Client // net/http Client to use for requests
	version       string       // v8
	url           string       // set of URLs passed initially to the client
	errorLog      Logger       // error log for critical messages
	infoLog       Logger       // information log for e.g. response times
	traceLog      Logger       // trace log for debugging
	token         string       // username for HTTP Basic Auth
	password      string       // password for HTTP Basic Auth
	maxRetries    uint
	sessionCookie string       //24 hour session cookie
	gzipEnabled   bool         // gzip compression enabled or disabled (default)
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*TogglClient) error

// An error is also returned when some configuration option is invalid
func NewClient(key string, options ...ClientOptionFunc) (*TogglClient, error) {
	// Set up the client
	c := &TogglClient{
		client:      http.DefaultClient,
		maxRetries:  DefaultMaxRetries,
		url:         DefaultUrl,
		version:     DefaultVersion,
		gzipEnabled: DefaultGzipEnabled,
		password:    DefaultAuthPassword,
		token:       key,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if len(c.token) < 10 {
		return nil, errors.New("Token required")
	}

	_, err := c.authenticate()

	if err != nil {
		return nil,err
	}

	if c.url != DefaultUrl {
		workspaceUrl = c.url + "/workspaces"
	}
	return c, nil
}

// SetHttpClient can be used to specify the http.Client to use when making
// HTTP requests to Elasticsearch.
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

// SetURL defines the URL endpoints of the Elasticsearch nodes. Notice that
// when sniffing is enabled, these URLs are used to initially sniff the
// cluster on startup.
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

func SetTraceLogger(l Logger) ClientOptionFunc {
	return func(c *TogglClient) error {
		c.traceLog = l
		return nil
	}
}

// SetURL defines the URL endpoints of the Elasticsearch nodes. Notice that
// when sniffing is enabled, these URLs are used to initially sniff the
// cluster on startup.
func SetToken(token string) ClientOptionFunc {
	return func(c *TogglClient) error {
		switch len(token) {
		case 0:
			return errors.New("Token required")
		default:
			c.token = token
		}
		return nil
	}
}

// SetMaxRetries sets the maximum number of retries before giving up when
// performing a HTTP request to Elasticsearch.
func SetMaxRetries(maxRetries uint) ClientOptionFunc {
	return func(c *TogglClient) error {
		if maxRetries < 0 {
			return errors.New("MaxRetries must be greater than or equal to 0")
		}
		c.maxRetries = maxRetries
		return nil
	}
}

// SetGzip enables or disables gzip compression (disabled by default).
func SetGzip(enabled bool) ClientOptionFunc {
	return func(c *TogglClient) error {
		c.gzipEnabled = enabled
		return nil
	}
}

func (c *TogglClient) String() string {
	return fmt.Sprintf("{token=%s}", c.token)
}

func (c *TogglClient) authenticate() ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s",c.url,"sessions"), nil)
	if err != nil {
		return nil, err
	}
	c.dumpRequest(req)
	req.SetBasicAuth(c.token, "api_token")
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
	return ioutil.ReadAll(resp.Body)
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

func (c *TogglClient) PostRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "POST", endpoint, body)
}

func (c *TogglClient) DeleteRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "DELETE", endpoint, body)
}

func (c *TogglClient) PutRequest(endpoint string, body []byte) ([]byte, error) {
	return request(c, "PUT", endpoint, body)
}

func (c *TogglClient) GetRequest(endpoint string) ([]byte, error) {
	return request(c, "GET", endpoint, nil)
}
