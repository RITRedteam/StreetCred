# StreetCred
Tool created for Red Team to test default credentials on SSH and WinRM and then execute scripts with those credentials before the password can be changed by Blue Team.

## Setup
The configuration file is located at `config/config.yml`. This is what will be looked for when the `-c` option is utilized.

If using the command line, just supply options as noted in the section below.

## Using Docker
If you plan on using Docker, a config file must be used.
```
docker build -t default .
docker run default
```
## Using the command line
```
go run main.go [options]
OR
go build main.go then ./main [options]
```
### Example command
`go run main.go -o output.txt -p password123 -u users.txt -b boxes.txt -s script.ps1`
### All Options
```
  -c
      If using a config file, use of this arg will set to bool value to true. Other arguments do not have to be provided.
  -b string
      Path to file containing list of boxes.
  -o string
      Output file name for successful responses. (default "output.txt")
  -p string
      Password to attempt on users and boxes.
  -s string
      Path to a script that should be executed on successful SSH/WinRM logon. If this option is not set, a script will not be executed.
  -u string
      Path to file containing list of users.
```
