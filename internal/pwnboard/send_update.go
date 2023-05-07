package pwnboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// This should be the URL/IP of the pwnboard instance that SendUpdate
// 	is sending the data to.
//var PWNBOARD string = "http://example.com" // Link to pwnboard [CHANGE ME]
var PWNBOARD, err = os.LookupEnv("PWNBOARD")

type Data struct {
	IPs  string `json:"ip"`   // Target IP address as a string
	Type string `json:"type"` // Describes what implant pwnboard is being updated from
}

// Handles error if obtaining environment variable fails
func CheckEnvVariable(pwnboard string, err bool) {
	if (len(strings.TrimSpace(pwnboard)) == 0) || (err != false) {
		os.Stderr.WriteString("ERROR: Failed to find environment variable")
		os.Stderr.WriteString(strconv.FormatBool(err) + "\n")
		return
	}
}

// Sends a post request with information about a target to pwnboard.
func SendUpdate(ip string, info string) {
	//use the Data struct to organize the data that will be sent to pwnboard
	data := Data{
		IPs:  ip,
		Type: info,
	}

	// Turn data struct into json
	mData, err := json.Marshal(data)
	if err != nil {
		os.Stderr.WriteString("ERROR: Failed to marshal data.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	// Send json data to pwnboard
	req, err := http.Post(PWNBOARD+"/generic", "application/json", bytes.NewBuffer(mData))
	if err != nil {
		os.Stderr.WriteString("ERROR: Failed to send a post request to pwnboard.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	// If anything is returned from pwnboard (usually nothing), print it to the terminal.
	var decoded map[string]interface{}
	json.NewDecoder(req.Body).Decode(&decoded)
	fmt.Println(decoded)
}
