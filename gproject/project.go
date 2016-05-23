package gproject

import (
	"encoding/json"
	"fmt"
	"github.com/dougEfresh/gtoggl"
)

// Toggl Project Definition
type Project struct {
	Id   uint64 `json:"id"`
	WId  uint64 `json:"wid"`
	CId  uint64 `json:"cid"`
	Name string `json:"name"`
}

type Projects []Project

const Endpoint = "/projects"

//Return a ProjectClient. An error is also returned when some configuration option is invalid
//    thc,err := gtoggl.NewClient("token")
//    pc,err := gproject.NewClient(tc)
func NewClient(thc *gtoggl.TogglHttpClient, options ...ProjectClientOptionFunc) (*ProjectClient, error) {
	tc := &ProjectClient{
		thc:             thc,
		listTransport:   defaultTransport,
		getTransport:    defaultTransport,
		updateTransport: defaultTransport,
		createTransport: defaultTransport,
		deleteTransport: defaultTransport,
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

type ProjectClient struct {
	thc             *gtoggl.TogglHttpClient
	endpoint        string
	listTransport   ClientLister
	getTransport    ClientGetter
	updateTransport ClientUpdater
	createTransport ClientCreater
	deleteTransport ClientDeleter
}

func (tc *ProjectClient) List() (Projects, error) {
	return tc.listTransport.List(tc)
}

func (tc *ProjectClient) Get(id uint64) (Project, error) {
	return tc.getTransport.Get(tc, id)
}

func (tc *ProjectClient) Create(c *Project) (Project, error) {
	return tc.createTransport.Create(tc, c)
}

func (tc *ProjectClient) Update(c *Project) (Project, error) {
	return tc.updateTransport.Update(tc, c)
}

func (tc *ProjectClient) Delete(id uint64) error {
	return tc.deleteTransport.Delete(tc, id)
}

type ClientLister interface {
	List(tc *ProjectClient) (Projects, error)
}
type ClientGetter interface {
	Get(tc *ProjectClient, id uint64) (Project, error)
}
type ClientUpdater interface {
	Update(tc *ProjectClient, c *Project) (Project, error)
}
type ClientCreater interface {
	Create(tc *ProjectClient, c *Project) (Project, error)
}
type ClientDeleter interface {
	Delete(tc *ProjectClient, id uint64) error
}

//Configures a Client.
/*
    func SetURL(url string) ToggleClientOptionFunc {
	return func(c *TogglClient) error {
	    c.Url = url
	}
    }
*/
type ProjectClientOptionFunc func(*ProjectClient) error
type projectTransport struct{}

var defaultTransport = &projectTransport{}

type projectResponse struct {
	Data Project `json:"data"`
}
type projectRequest struct {
	Data Project `json:"data"`
}
type projectUpdateRequest struct {
	Project Project `json:"project"`
}

//GET https://www.toggl.com/api/v8/projects/1239455
func (cl *projectTransport) Get(tc *ProjectClient, id uint64) (Project, error) {
	body, err := tc.thc.GetRequest(fmt.Sprintf("%s/%d", tc.endpoint, id))
	if err != nil {
		return Project{}, err
	}

	var aux projectResponse
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

//DELETE https://www.toggl.com/api/v8/projects/1239455
func (cl *projectTransport) Delete(tc *ProjectClient, id uint64) error {
	_, err := tc.thc.DeleteRequest(fmt.Sprintf("%s/%d", tc.endpoint, id), nil)
	return err
}

//GET https://www.toggl.com/api/v8/projects/1239455
func (cl *projectTransport) List(tc *ProjectClient) (Projects, error) {
	body, err := tc.thc.GetRequest(tc.endpoint)
	var Clients []Project
	if err != nil {
		return Clients, err
	}
	err = json.Unmarshal(body, &Clients)
	return Clients, err
}

//PUT https://www.toggl.com/api/v8/projects/1239455
func (cl *projectTransport) Update(tc *ProjectClient, c *Project) (Project, error) {
	put := projectUpdateRequest{Project: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return Project{}, err
	}
	response, err := tc.thc.PutRequest(fmt.Sprintf("%s/%d", tc.endpoint, c.Id), body)
	if err != nil {
		return Project{}, err
	}
	var aux projectResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

//POST https://www.toggl.com/api/v8/projects
func (cl *projectTransport) Create(tc *ProjectClient, c *Project) (Project, error) {
	put := projectUpdateRequest{Project: *c}
	body, err := json.Marshal(put)
	if err != nil {
		return Project{}, err
	}
	response, err := tc.thc.PostRequest(tc.endpoint, body)
	if err != nil {
		return Project{}, err
	}
	var aux projectResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}
