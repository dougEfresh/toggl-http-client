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
		thc:             thc,
		listTransport:   defaultClientTransport,
		getTransport:    defaultClientTransport,
		updateTransport: defaultClientTransport,
		createTransport: defaultClientTransport,
		deleteTransport: defaultClientTransport,
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
	thc             *gtoggl.TogglHttpClient
	endpoint        string
	listTransport   ClientLister
	getTransport    ClientGetter
	updateTransport ClientUpdater
	createTransport ClientCreater
	deleteTransport ClientDeleter
}

func (tc *TogglClient) List() (Clients, error) {
	return tc.listTransport.List(tc)
}

func (tc *TogglClient) Get(id uint64) (Client, error) {
	return tc.getTransport.Get(tc, id)
}

func (tc *TogglClient) Create(c *Client) (Client, error) {
	return tc.createTransport.Create(tc, c)
}

func (tc *TogglClient) Update(c *Client) (Client, error) {
	return tc.updateTransport.Update(tc, c)
}

func (tc *TogglClient) Delete(id uint64) error {
	return tc.deleteTransport.Delete(tc, id)
}

type ClientLister interface {
	List(tc *TogglClient) (Clients, error)
}
type ClientGetter interface {
	Get(tc *TogglClient, id uint64) (Client, error)
}
type ClientUpdater interface {
	Update(tc *TogglClient, c *Client) (Client, error)
}
type ClientCreater interface {
	Create(tc *TogglClient, c *Client) (Client, error)
}
type ClientDeleter interface {
	Delete(tc *TogglClient, id uint64) error
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
type clientTransport struct{}

var defaultClientTransport = &clientTransport{}

type clientResponse struct {
	Data Client `json:"data"`
}
type clientRequest struct {
	Data Client `json:"data"`
}
type clientCreateRequest struct {
	Client Client `json:"client"`
}

//GET https://www.toggl.com/api/v8/clients/1239455
func (cl *clientTransport) Get(tc *TogglClient, id uint64) (Client, error) {
	body, err := tc.thc.GetRequest(fmt.Sprintf("%s/%d", tc.endpoint, id))
	if err != nil {
		return Client{}, err
	}

	var aux clientResponse
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

//DELETE https://www.toggl.com/api/v8/clients/1239455
func (cl *clientTransport) Delete(tc *TogglClient, id uint64) error {
	_, err := tc.thc.DeleteRequest(fmt.Sprintf("%s/%d", tc.endpoint, id), nil)
	return err
}

//GET https://www.toggl.com/api/v8/clients/1239455
func (cl *clientTransport) List(tc *TogglClient) (Clients, error) {
	body, err := tc.thc.GetRequest(tc.endpoint)
	var Clients []Client
	if err != nil {
		return Clients, err
	}
	err = json.Unmarshal(body, &Clients)
	return Clients, err
}

//PUT https://www.toggl.com/api/v8/clients/1239455
func (cl *clientTransport) Update(tc *TogglClient, c *Client) (Client, error) {
	put := clientCreateRequest{Client: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return Client{}, err
	}
	response, err := tc.thc.PutRequest(fmt.Sprintf("%s/%d", tc.endpoint, c.Id), body)
	if err != nil {
		return Client{}, err
	}
	var aux clientResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

//POST https://www.toggl.com/api/v8/clients
func (cl *clientTransport) Create(tc *TogglClient, c *Client) (Client, error) {
	put := clientCreateRequest{Client: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return Client{}, err
	}
	response, err := tc.thc.PostRequest(tc.endpoint, body)
	if err != nil {
		return Client{}, err
	}
	var aux clientResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}
