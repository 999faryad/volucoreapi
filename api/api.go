package api

import (
	"VoluCore/api/crontab"
	"VoluCore/api/firewall"
	"VoluCore/api/getload"
	"net/http"
)

func OpenWebserver() error {
	http.HandleFunc("/getload", getload.GetLoad)
	http.HandleFunc("/addfwrule", firewall.AddFWRule)
	http.HandleFunc("/delfwrule", firewall.DelFWRule)
	http.HandleFunc("/getfwrules", firewall.GetFWRules)
	http.HandleFunc("/addcron", crontab.AddCronJob)
	err := http.ListenAndServe("0.0.0.0:9890", nil)
	if err != nil {
		return err
	}
	return nil
}
