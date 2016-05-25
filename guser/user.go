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
		thc: thc,
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
}

func (c *UserClient) Get(realatedData bool) (*User, error) {
	return userResponse(c.thc.GetRequest(c.endpoint))
}

func (c *UserClient) Create(email, password, timezone string) (*User, error) {
	up := &UserCreate{Password: password, Email: email, Timezone: timezone, CreatedWith: "gtoggl"}
	put := createRequest{User: up}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return userResponse(c.thc.PostRequest(c.signupEndpoint, body))
}

func (c *UserClient) Update(u *User) (*User, error) {
	up := &UserUpdate{FullName: u.FullName, Email: u.Email}
	put := updateRequest{User: up}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return userResponse(c.thc.PutRequest(c.endpoint, body))
}

func (c *UserClient) ResetToken() (string, error) {
	response, err := c.thc.PostRequest(c.resetEndpoint, nil)
	if err != nil {
		return "", err
	}
	var aux string
	err = json.Unmarshal(*response, &aux)
	return aux, nil

}

func userResponse(response *json.RawMessage, error error) (*User, error) {
	if error != nil {
		return nil, error
	}
	var tResp gtoggl.TogglResponse
	err := json.Unmarshal(*response, &tResp)
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal(*tResp.Data, &u)
	if err != nil {
		return nil, err
	}
	return &u, err
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

type request struct {
	Data User `json:"data"`
}
type updateRequest struct {
	User *UserUpdate `json:"user"`
}
type createRequest struct {
	User *UserCreate `json:"user"`
}
