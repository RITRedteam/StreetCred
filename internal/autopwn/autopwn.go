package autopwn

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

// Will attempt to execute a script located at scriptPath on the target host using
//	provided user and password through SSH.
func SSHAutopwn(host, user, password, scriptPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Set up SSH connection config
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout:         30 * time.Second,
	}

	// Attempt SSH connection/login
	conn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Could not log into SSH with procided host, user, and password.")
		os.Stderr.WriteString(err.Error())
		return
	}
	fmt.Println("autopwn: Successful SSH connection @", host)

	// Close the connection when the rest of the function is done running
	defer conn.Close()

	// TODO: Read script from scriptPath
	// TODO: Execute script through SSH on host

}

// Will attempt to execute a script located at scriptPath on the target host using
//	provided user and password through WinRM.
func WinRMAutopwn(host, user, password, scriptPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Split host and port number to be used when creating an endpoint
	splitHost := strings.Split(host, ":")
	port, err := strconv.Atoi(splitHost[1])
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to convert port number into int.")
		os.Stderr.WriteString(err.Error())
		return
	}

	// Create an endpoint and setup the WinRM connection
	endpoint := winrm.NewEndpoint(splitHost[0], port, false, true, nil, nil, nil, 30*time.Second)
	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	// Attempt to create WinRM client
	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to create WinRM client")
		os.Stderr.WriteString(err.Error())
		return
	}

	// Attempt to execute a basic powershell command through WinRM
	cmd := winrm.Powershell("ipconfig")
	_, err = client.Run(cmd, os.Stdout, os.Stderr)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to execute command through WinRM.")
		os.Stderr.WriteString(err.Error())
		return
	}

	fmt.Println("Successful WinRM connection @", host)

	// TODO: Read script from scriptPath
	// TODO: Execute script through WinRM on host
}
