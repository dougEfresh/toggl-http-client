package gtoggl

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	Id       uint64 `json:"id"`
	WId      uint64 `json:"wid"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type Clients []Client

type ClientLister interface {
	List(tc *TogglClient) (Clients, error)
}
type ClientGetter interface {
	Get(tc *TogglClient, id uint64) (Client, error)
}
type ClientUpdater interface {
	Update(tc *TogglClient, c Client) (Client, error)
}
type ClientCreater interface {
	Create(tc *TogglClient, c Client) (Client, error)
}

// ClientOptionFunc is a function that configures a Client.
// It is used in NewTogglClient.
type ToggleClientOptionFunc func(*TogglClient) error
type clientTransport struct{}

var defaultClientTransport = &clientTransport{}
type clientResponse struct {
	Data Client `json:"data"`
}
type clientRequest struct {
	Data Client `json:"data"`
}

func (cl *clientTransport) Get(wsc *TogglClient, id uint64) (Client, error) {
	body, err := wsc.tc.GetRequest(fmt.Sprintf("%s/%d", wsc.clientEndpoint, id))
	if err != nil {
		return Client{}, err
	}

	var aux clientResponse
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

func (cl *clientTransport) List(tc *TogglClient) (Clients, error) {
	body, err := tc.tc.GetRequest(tc.clientEndpoint)
	var Clients []Client
	if err != nil {
		return Clients, err
	}
	err = json.Unmarshal(body, &Clients)
	return Clients, err
}

func (cl *clientTransport) Update(tc *TogglClient, c Client) (Client, error) {
	put := clientRequest{Data: c}
	body, err := json.Marshal(put)
	if err != nil {
		return Client{}, err
	}
	response, err := tc.tc.PutRequest(fmt.Sprintf("%s/%d", tc.clientEndpoint, c.Id), body)
	if err != nil {
		return Client{}, err
	}
	var aux clientResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

type TogglClient struct {
	tc              *TogglHttpClient
	clientEndpoint  string
	listTransport   ClientLister
	getTransport    ClientGetter
	updateTransport ClientUpdater
	createTransport ClientCreater
}

func NewTogglClient(thc *TogglHttpClient, options...ToggleClientOptionFunc) (*TogglClient,error) {
	tc := &TogglClient{
		tc:              thc,
		listTransport:   defaultClientTransport,
		getTransport:    defaultClientTransport,
		updateTransport: defaultClientTransport,
	}
	// Run the options on it

	for _, option := range options {
		if err := option(tc); err != nil {
			return nil, err
		}
	}
	tc.clientEndpoint = thc.Url + "/clients"
	return tc, nil
}

func (tc *TogglClient) List() (Clients,error) {
	return tc.listTransport.List(tc)
}
func (tc *TogglClient) Get(id uint64) (Client,error) {
	return tc.getTransport.Get(tc,id)
}
