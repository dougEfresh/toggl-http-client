package gtoggl

type Client struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}
