package gtimeentry

import (
	"encoding/json"
	"fmt"
	"github.com/dougEfresh/gtoggl"
	"github.com/dougEfresh/gtoggl/gproject"
	"github.com/dougEfresh/gtoggl/gworkspace"
	"time"
)

type TimeEntry struct {
	Id          uint64                `json:"id,omitempty"`
	Description string                `json:"description"`
	Project     *gproject.Project     `json:"project"`
	Start       time.Time             `json:"start"`
	Stop        time.Time             `json:"stop"`
	Duration    uint64                `json:"duration"`
	Billable    bool                  `json:"billable"`
	Workspace   *gworkspace.Workspace `json:"workspace"`
	Tags        []string              `json:"tags"`
	Pid         uint64                `json:"pid"`
	Wid         uint64                `json:"wid"`
	Tid         uint64                `json:"tid"`
	CreatedWith string                `json:"created_with,omitempty" `
}

type TimeEntries []TimeEntry

const Endpoint = "/time_entries"
const EndpointCurrent = Endpoint + "/current"
const EndpointStart = Endpoint + "/start"

//Return a UserClient. An error is also returned when some configuration option is invalid
//    thc,err := gtoggl.NewClient("token")
//    uc,err := guser.NewClient(thc)
func NewClient(thc *gtoggl.TogglHttpClient, options ...ClientOptionFunc) (*TimeEntryClient, error) {
	tc := &TimeEntryClient{
		thc: thc,
	}
	// Run the options on it
	for _, option := range options {
		if err := option(tc); err != nil {
			return nil, err
		}
	}
	tc.endpoint = thc.Url + Endpoint
	tc.currentEndpoint = thc.Url + EndpointCurrent
	tc.startEndpoint = thc.Url + EndpointStart
	return tc, nil
}

type TimeEntryClient struct {
	thc             *gtoggl.TogglHttpClient
	endpoint        string
	startEndpoint   string
	currentEndpoint string
}

func (c *TimeEntryClient) Get(tid uint64) (*TimeEntry, error) {
	return timeEntryResponse(c.thc.GetRequest(fmt.Sprintf("%s/%d", c.endpoint, tid)))
}

func (tc *TimeEntryClient) Delete(id uint64) error {
	_, err := tc.thc.DeleteRequest(fmt.Sprintf("%s/%d", tc.endpoint, id), nil)
	return err
}

func (c *TimeEntryClient) List() (TimeEntries, error) {
	body, err := c.thc.GetRequest(c.endpoint)
	var te TimeEntries
	if err != nil {
		return te, err
	}
	err = json.Unmarshal(*body, &te)
	return te, err
}

func (c *TimeEntryClient) Create(t *TimeEntry) (*TimeEntry, error) {
	if len(t.CreatedWith) < 0 {
		t.CreatedWith = gtoggl.TogglCreator
	}
	up := createRequest{TimeEntry: t}
	body, err := json.Marshal(up)
	if err != nil {
		return nil, err
	}
	return timeEntryResponse(c.thc.PutRequest(c.endpoint, body))
}

func (c *TimeEntryClient) Update(t *TimeEntry) (*TimeEntry, error) {
	up := updateRequest{TimeEntry: t}
	body, err := json.Marshal(up)
	if err != nil {
		return nil, err
	}
	return timeEntryResponse(c.thc.PutRequest(fmt.Sprintf("%s/%d", c.endpoint, t.Id), body))
}

func timeEntryResponse(response *json.RawMessage, error error) (*TimeEntry, error) {
	if error != nil {
		return nil, error
	}
	var tResp gtoggl.TogglResponse
	err := json.Unmarshal(*response, &tResp)
	if err != nil {
		return nil, err
	}
	var t TimeEntry
	err = json.Unmarshal(*tResp.Data, &t)
	if err != nil {
		return nil, err
	}
	return &t, err
}

//Configures a Client.
/*
    func SetURL(url string) ToggleClientOptionFunc {
	return func(c *TogglClient) error {
	    c.Url = url
	}
    }
*/
type ClientOptionFunc func(*TimeEntryClient) error

type updateRequest struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}
type createRequest struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}
