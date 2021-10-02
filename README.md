# bruhdotzip
Tool created for Red Team to test default credentials on SSH and WinRM and then execute scripts with those credentials before the password can be changed by Blue Team.

## Use
`go run main.go [options]`
OR
`go build main.go` then `./main [options]`

```
OPTIONS:
  -b string
        Path to file containing list of boxes. (shorthand)
  -boxPath string
        Path to file containing list of boxes.
  -o string
        Output file name for successful responses. (shorthand) (default "output.txt")
  -output string
        Output file name for successful responses. (default "output.txt")
  -p string
        Password to attempt on users and boxes. (shorthand)
  -password string
        Password to attempt on users and boxes.
  -s string
        Path to a script that should be executed on successful SSH/WinRM logon. If this option is not set, a script will not be executed. (shorthand)
  -script string
        Path to a script that should be executed on successful SSH/WinRM logon. If this option is not set, a script will not be executed.
  -u string
        Path to file containing list of users. (shorthand)
  -userPath string
        Path to file containing list of users.
```

### Example command:
  `go run main.go -o output.txt -p password123 -u users.txt -b boxes.txt -s script.ps1`
