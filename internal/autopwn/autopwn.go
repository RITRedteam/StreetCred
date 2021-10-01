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
	fmt.Println("autopwn: Successful SSH connection @", host)
	defer conn.Close()

	// Read contents of script file and save to a string to be sent over through ssh
	fileString, err := files.ReadString(scriptPath)
	fmt.Println(fileString)

	// Create a new file in the tmp dir on the remote host containing the contents of the script
	sshSessionExec(conn, "echo \""+fileString+"\"> /tmp/output.sh", host)
	// Execute the script on the remote host
	sshSessionExec(conn, "sh /tmp/output.sh", host)
	// Delete the script on the remote host such that the same path can be used again (and less trail)
	sshSessionExec(conn, "rm /tmp/output.sh", host)

}

func sshSessionExec(conn *ssh.Client, cmd string, host string) error {
	session, err := conn.NewSession()
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Could not create an SSH session.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
	fmt.Println("autopwn: Successful SSH session creation @", host)

	err = session.Run(cmd)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to execute command [ " + cmd + "] on the remote host.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
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

	// Attempt to execute a powershell command through WinRM
	fileString, err := files.ReadString(scriptPath)
	//if scriptPath[len(scriptPath)-3:] == "bat"
	fmt.Println(fileString)
	cmd := winrm.Powershell("echo \"" + fileString + "\" > 'C:/Windows/Temp/output.bat'")
	_, err = client.Run(cmd, os.Stdout, os.Stderr)
	if err != nil {
		os.Stderr.WriteString("ERROR(autopwn): Failed to write script to a file on the remote host through WinRM.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
	cmd = winrm.Powershell("'C:/Windows/Temp/output.bat'")
	_, err = client.Run(cmd, os.Stdout, os.Stderr)

	// TODO: Read script from scriptPath
	// TODO: Execute script through WinRM on host
}
