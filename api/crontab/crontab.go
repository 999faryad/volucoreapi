package crontab

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type CronJobRequest struct {
	Minute  string `json:"minute"`
	Hour    string `json:"hour"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	Weekday string `json:"weekday"`
	Command string `json:"command"`
}

func (c *CronJobRequest) ToCronCommand() string {
	return fmt.Sprintf("%s %s %s %s %s %s", c.Minute, c.Hour, c.Day, c.Month, c.Weekday, c.Command)
}

func Respond(w http.ResponseWriter, responseError bool, responseMessage string, responseStatus *int) {
	responseStruct := &Response{
		Error:   responseError,
		Message: responseMessage,
	}
	response, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("FATAL Error while Responding. Please open an Issue and specify following: %v", err.Error())
	}

	w.WriteHeader(*responseStatus)
	fmt.Fprintf(w, string(response))
}
