package firewall

import (
	"fmt"
	ipt "github.com/coreos/go-iptables/iptables"
	"log"
	"net/http"
)

func GetFWRules(writer http.ResponseWriter, request *http.Request) {
	response := &RuleResponse{}
	chain := INPUTCHAIN
	table := FILTERTABLE

	iptables, err := ipt.New()
	if err != nil {
		currentError := fmt.Sprintf("Unable to initialize new IPTables instance. \nDetailed Error:\n%v", err.Error())
		log.Printf(currentError)
		response = &RuleResponse{
			Rules: []string{currentError},
		}
		RespondJSON(writer, *response)
		return
	}

	fwRules, err := iptables.List(table, chain)
	if err != nil {
		currentError := fmt.Sprintf("Error occured while trying to access chain: [%v] on table: [%v].\nDetailed Error:\n%v", chain, table, err.Error())
		log.Printf(currentError)
		response = &RuleResponse{
			Rules: []string{currentError},
		}
		RespondJSON(writer, *response)
		return
	}
	response = &RuleResponse{
		Rules: fwRules,
	}
	RespondJSON(writer, *response)
}
