package firewall

import (
	"fmt"
	ipt "github.com/coreos/go-iptables/iptables"
	"log"
	"net/http"
	"os/exec"
)

func AddFWRule(writer http.ResponseWriter, request *http.Request) {
	respJson := &Response{}
	table := FILTERTABLE

	iptables, err := ipt.New()
	if err != nil {
		currentError := fmt.Sprintf("Unable to initialize new IPTables instance. \nDetailed Error:\n%v", err.Error())
		log.Printf(currentError)
		respJson = &Response{
			Error:   true,
			Message: currentError,
		}
		RespondJSON(writer, *respJson)
		return
	}

	data := GetRequestData(request)
	rule := []string{"-s", data.SAddr, "-p", data.Protocol, "--dport", data.DPort, "-j", data.Policy}

	if err := iptables.AppendUnique(table, INPUTCHAIN, rule...); err != nil {
		currentError := fmt.Sprintf("Unable to add new Firewall Rule. \nDetailed Error:\n%v", err.Error())
		log.Printf(currentError)
		respJson = &Response{
			Error:   true,
			Message: currentError,
		}
		RespondJSON(writer, *respJson)
		return
	}

	if err := exec.Command("iptables-save", "-c").Run(); err != nil {
		currentError := fmt.Sprintf("Error saving IPTables Rules to IPTables Config. \nDetailed Error:\n%v", err.Error())
		log.Printf(currentError)
		respJson = &Response{
			Error:   true,
			Message: currentError,
		}
		RespondJSON(writer, *respJson)
		return
	}

	respJson = &Response{
		Error:   false,
		Message: "Successfully added Firewall Rule.",
	}
	RespondJSON(writer, *respJson)
}
