package firewall

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	INPUTCHAIN   string = "INPUT"
	OUTPUTCHAIN  string = "OUTPUT"
	FORWARDCHAIN string = "FORWARD"

	FILTERTABLE string = "filter"
)

type RuleResponse struct {
	Rules []string `json:"rules"`
}

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

type ResponseStructs interface {
	RuleResponse | Response
}

func GetRequestData(request *http.Request) *Rule {
	rule := &Rule{}
	err := json.NewDecoder(request.Body).Decode(&rule)
	if err != nil {
		return rule
	}
	return rule
}

func RespondJSON[T ResponseStructs](writer http.ResponseWriter, r T) {

	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Fprintf(writer, "Unable to Marshal JSON!")
	}
	fmt.Fprintf(writer, string(response))

}
