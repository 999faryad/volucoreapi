package crontab

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func AddCronJob(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	if r.Method != http.MethodPost {
		statusCode = http.StatusMethodNotAllowed
		Respond(w, true, "Invalid Method for this API Call.", &statusCode)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		statusCode = http.StatusBadRequest
		Respond(w, true, "Could not read request Body", &statusCode)
		return
	}
	bodyData := &CronJobRequest{}
	err = json.Unmarshal(body, &bodyData)
	if body == nil || err != nil {
		statusCode = http.StatusBadRequest
		fmt.Println(err.Error())
		Respond(w, true, "Please specify the Parameters to add a cronjob", &statusCode)
		return
	}
	remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		statusCode = http.StatusInternalServerError
		Respond(w, true, "Error while getting Caller Address", &statusCode)
		return
	}

	if remoteAddr != "127.0.0.1" && remoteAddr != "localhost" {
		statusCode = http.StatusForbidden
		Respond(w, true, "Security Block!", &statusCode)
		return
		// Maybe add some error Reporting to Customer, because it seems, that the API is accessible from Outside.
	}

	createSystemCronJob(w, bodyData.ToCronCommand())

	Respond(w, true, "Successfully added Cronjob", nil)
	return
}

func createSystemCronJob(w http.ResponseWriter, schedule string) {
	tempFile, err := os.CreateTemp("", "temp-cron")
	if err != nil {
		statusCode := http.StatusInternalServerError
		Respond(w, true, "Error while creating Temporary File", &statusCode)
		return
	}
	defer os.Remove(tempFile.Name())

	cmd := exec.Command("crontab", "-l")
	cmdOutput, err := cmd.Output()
	if err != nil {
		statusCode := http.StatusInternalServerError
		Respond(w, true, "Error while executing crontab -l", &statusCode)
		return
	}

	newCronJob := fmt.Sprintf("%s\n", schedule)
	err = os.WriteFile(tempFile.Name(), []byte(string(cmdOutput)+newCronJob), 0644)
	if err != nil {
		statusCode := http.StatusInternalServerError
		Respond(w, true, fmt.Sprintf("Error while writing to file %v", tempFile.Name()), &statusCode)
		return
	}

	cmd = exec.Command("crontab", tempFile.Name())
	err = cmd.Run()
	if err != nil {
		statusCode := http.StatusInternalServerError
		Respond(w, true, fmt.Sprintf("Error while executing command crontab %v", tempFile.Name()), &statusCode)
		return
	}
}
