package guser

import (
	"encoding/json"
	"github.com/dougEfresh/gtoggl"
)

// Toggl User Definition
type User struct {
	Id       uint64 `json:"id"`
	ApiToken string `json:"api_token"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}
type UserUpdate struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}
type UserCreate struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Timezone    string `json:"timezone"`
	CreatedWith string `json:"created_with"`
}

type Users []User

const Endpoint = "/me"
const SignupEndpoint = "/signups"
const ResetEndpoint = "/reset_token"
const MeWithRelatedData = "/me?with_related_data=true"

//Return a UserClient. An error is also returned when some configuration option is invalid
//    thc,err := gtoggl.NewClient("token")
//    uc,err := guser.NewClient(thc)
func NewClient(thc *gtoggl.TogglHttpClient, options ...ClientOptionFunc) (*UserClient, error) {
	tc := &UserClient{
		thc:             thc,
		getTransport:    defaultTransport,
		updateTransport: defaultTransport,
		createTransport: defaultTransport,
		resetTransport:  defaultTransport,
	}
	// Run the options on it
	for _, option := range options {
		if err := option(tc); err != nil {
			return nil, err
		}
	}
	tc.endpoint = thc.Url + Endpoint
	tc.signupEndpoint = thc.Url + SignupEndpoint
	tc.resetEndpoint = thc.Url + ResetEndpoint
	tc.relatedEndpoint = thc.Url + MeWithRelatedData
	return tc, nil
}

type UserClient struct {
	thc             *gtoggl.TogglHttpClient
	endpoint        string
	resetEndpoint   string
	signupEndpoint  string
	relatedEndpoint string
	getTransport    UserGetter
	updateTransport UserUpdater
	createTransport UserCreator
	resetTransport  UserResetter
}

func (c *UserClient) Get(realatedData bool) (User, error) {
	return c.getTransport.Get(c, realatedData)
}

func (c *UserClient) Create(email, password, timezone string) (User, error) {
	return c.createTransport.Create(c, email, password, timezone)
}

func (c *UserClient) Update(u *User) (User, error) {
	return c.updateTransport.Update(c, u)
}

func (c *UserClient) ResetToken() (string, error) {
	return c.resetTransport.ResetToken(c)
}

type UserGetter interface {
	Get(tc *UserClient, relatedData bool) (User, error)
}
type UserUpdater interface {
	Update(tc *UserClient, c *User) (User, error)
}
type UserCreator interface {
	Create(tc *UserClient, email, password, timezone string) (User, error)
}
type UserResetter interface {
	ResetToken(tc *UserClient) (string, error)
}

//Configures a Client.
/*
    func SetURL(url string) ToggleClientOptionFunc {
	return func(c *TogglClient) error {
	    c.Url = url
	}
    }
*/
type ClientOptionFunc func(*UserClient) error
type transport struct{}

var defaultTransport = &transport{}

type tResponse struct {
	Data User `json:"data"`
}
type request struct {
	Data User `json:"data"`
}
type updateRequest struct {
	User *UserUpdate `json:"user"`
}
type createRequest struct {
	User *UserCreate `json:"user"`
}

//GET https://www.toggl.com/api/v8/me
func (ct *transport) Get(tc *UserClient, relatedData bool) (User, error) {
	body, err := tc.thc.GetRequest(tc.endpoint)
	if err != nil {
		return User{}, err
	}

	var aux tResponse
	err = json.Unmarshal(body, &aux)
	return aux.Data, err
}

//PUT https://www.toggl.com/api/v8/me/1239455
func (ct *transport) Update(tc *UserClient, c *User) (User, error) {
	up := &UserUpdate{FullName: c.FullName, Email: c.Email}
	put := updateRequest{User: up}
	body, err := json.Marshal(put)
	if err != nil {
		return User{}, err
	}
	response, err := tc.thc.PutRequest(tc.endpoint, body)
	if err != nil {
		return User{}, err
	}
	var aux tResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

//POST https://www.toggl.com/api/v8/signups
func (ct *transport) Create(tc *UserClient, email, password, timezone string) (User, error) {
	up := &UserCreate{Password: password, Email: email, Timezone: timezone, CreatedWith: "gtoggl"}
	put := createRequest{User: up}
	body, err := json.Marshal(put)
	if err != nil {
		return User{}, err
	}
	response, err := tc.thc.PostRequest(tc.signupEndpoint, body)
	if err != nil {
		return User{}, err
	}
	var aux tResponse
	err = json.Unmarshal(response, &aux)
	return aux.Data, err
}

//POST https://www.toggl.com/api/v8/reset_token
func (ct *transport) ResetToken(tc *UserClient) (string, error) {
	response, err := tc.thc.PostRequest(tc.resetEndpoint, nil)
	if err != nil {
		return "", err
	}
	var aux string
	err = json.Unmarshal(response, &aux)
	return aux, nil
}
