package api

import (
	"VoluCore/api/firewall"
	"VoluCore/api/getload"
	"net/http"
)

func OpenWebserver() {
	http.HandleFunc("/getload", getload.GetLoad)
	http.HandleFunc("/addfwrule", firewall.AddFWRule)
	http.HandleFunc("/delfwrule", firewall.DelFWRule)
	http.HandleFunc("/getfwrules", firewall.GetFWRules)
	err := http.ListenAndServe("0.0.0.0:9890", nil)
	if err != nil {
		return
	}
}
