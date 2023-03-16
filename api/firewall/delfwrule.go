package firewall

import (
	"fmt"
	ipt "github.com/coreos/go-iptables/iptables"
	"log"
	"net/http"
	"os/exec"
)

func DelFWRule(writer http.ResponseWriter, request *http.Request) {
	respJson := &Response{}
	iptables, err := ipt.New()
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return
	}
	table := "filter"
	data := GetRequestData(request)
	rule := []string{"-s", data.SAddr, "-p", data.Protocol, "--dport", data.DPort, "-j", data.Policy}

	if err := iptables.Delete(table, INPUTCHAIN, rule...); err != nil {
		log.Print(err.Error())
		respJson = &Response{
			Error:   true,
			Message: err.Error(),
		}
	} else {
		respJson = &Response{
			Error:   false,
			Message: "Successfully removed Firewall Rule.",
		}
	}

	if err := exec.Command("iptables-save", "-c").Run(); err != nil {
		fmt.Printf("Error saving IPTables Config: %v", err)
		respJson = &Response{
			Error:   false,
			Message: fmt.Sprintf("Could not Save iptables for restart. NOT PERSISTENT, ERROR: %v", err.Error()),
		}
		return
	} else {
		fmt.Println("Successfully saved IPTables Config File")
	}

	fmt.Fprintf(writer, respJson.GetJsonResponse())

}
