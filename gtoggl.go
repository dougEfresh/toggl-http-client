package gtoggl

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var _client_path = "/clients"

type Gtoggl struct {
	Token string
	Url   string
}

var DefaultGtoggl = &Gtoggl{Url: "https://www.toggl.com/api/v8"}

func (g *Gtoggl) clients() []Client {
	req, _ := http.NewRequest("GET", g.Url+_client_path, nil)
	req.SetBasicAuth(g.Token, "api_token")
	hclient := &http.Client{}
	resp, err := hclient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var clients []Client
	json.Unmarshal([]byte(body), &clients)
	return clients
}
