package main

import (
	"VoluCore/api"
	"log"
)

func main() {
	err := api.OpenWebserver()
	if err != nil {
		log.Fatalf("An Error Occured while Opening Webserver\nError:\n%v", err.Error())
	}
}
