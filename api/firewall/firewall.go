package firewall

import (
	"encoding/json"
	"net/http"
)

const (
	INPUTCHAIN string = "INPUT"
)

type Rule struct {
	SAddr    string `json:"saddr"`
	DPort    string `json:"dport"`
	Policy   string `json:"policy"`
	Protocol string `json:"protocol"`
}

type Response struct {
	Error   bool
	Message string
}

func GetRequestData(request *http.Request) *Rule {
	rule := &Rule{}
	err := json.NewDecoder(request.Body).Decode(&rule)
	if err != nil {
		return rule
	}
	return rule
}

func (r *Response) GetJsonResponse() string {

	response, err := json.Marshal(&r)
	if err != nil {
		return ""
	}
	return string(response)

}
