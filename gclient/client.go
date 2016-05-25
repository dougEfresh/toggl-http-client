package gclient

import (
	"encoding/json"
	"fmt"
	"github.com/dougEfresh/gtoggl"
)

// Toggl Client Definition
type Client struct {
	Id       uint64 `json:"id"`
	WId      uint64 `json:"wid"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type Clients []Client

const Endpoint = "/clients"

//Return a Toggl Client. An error is also returned when some configuration option is invalid
//    thc,err := gtoggl.NewClient("token")
//    tc,err := gclient.NewClient(tc)
func NewClient(thc *gtoggl.TogglHttpClient, options ...ToggleClientOptionFunc) (*TogglClient, error) {
	tc := &TogglClient{
		thc: thc,
	}
	// Run the options on it
	for _, option := range options {
		if err := option(tc); err != nil {
			return nil, err
		}
	}
	tc.endpoint = thc.Url + Endpoint
	return tc, nil
}

type TogglClient struct {
	thc      *gtoggl.TogglHttpClient
	endpoint string
}

func (tc *TogglClient) List() (Clients, error) {
	body, err := tc.thc.GetRequest(tc.endpoint)
	var clients Clients
	if err != nil {
		return clients, err
	}
	err = json.Unmarshal(*body, &clients)
	return clients, err
}

func (tc *TogglClient) Get(id uint64) (*Client, error) {
	return clientResponse(tc.thc.GetRequest(fmt.Sprintf("%s/%d", tc.endpoint, id)))
}

func (tc *TogglClient) Create(c *Client) (*Client, error) {
	put := clientCreateRequest{Client: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return clientResponse(tc.thc.PostRequest(tc.endpoint, body))
}

func (tc *TogglClient) Update(c *Client) (*Client, error) {
	put := clientCreateRequest{Client: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return clientResponse(tc.thc.PutRequest(fmt.Sprintf("%s/%d", tc.endpoint, c.Id), body))
}

func (tc *TogglClient) Delete(id uint64) error {
	_, err := tc.thc.DeleteRequest(fmt.Sprintf("%s/%d", tc.endpoint, id), nil)
	return err
}

//Configures a Client.
/*
    func SetURL(url string) ToggleClientOptionFunc {
	return func(c *TogglClient) error {
	    c.Url = url
	}
    }
*/
type ToggleClientOptionFunc func(*TogglClient) error

type clientCreateRequest struct {
	Client Client `json:"client"`
}

func clientResponse(response *json.RawMessage, error error) (*Client, error) {
	if error != nil {
		return nil, error
	}
	var tResp gtoggl.TogglResponse
	err := json.Unmarshal(*response, &tResp)
	if err != nil {
		return nil, err
	}
	var cl Client
	err = json.Unmarshal(*tResp.Data, &cl)
	if err != nil {
		return nil, err
	}
	return &cl, err
}
