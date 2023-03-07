package api

import (
	"VoluCore/api/getload"
	"net/http"
)

func OpenWebserver() {
	http.HandleFunc("/getload", getload.GetLoad)
	err := http.ListenAndServe("127.0.0.1:9890", nil)
	if err != nil {
		return
	}
}
