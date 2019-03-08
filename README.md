# README

**SUMMARY**: JBlastor is a command line utility written in go for parsing JSON
files then performing a POST request to the target REST API endpoint.


## USAGE

```
➜  jblastor git:(master) ✗ ./jblastor
jblastor: error: required flag --files not provided, try --help
➜  jblastor git:(master) ✗

➜  jblastor git:(master) ✗ ./jblastor --help
usage: jblastor --files=FILES --endpoint=ENDPOINT [<flags>]

Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:8088/save'

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -f, --files=FILES        Path to file or directory of file(s) to parse and POST.
  -e, --endpoint=ENDPOINT  REST API endpoint to send request to.
  -u, --apiuser=APIUSER    API User account permitted to do POST requests.
  -p, --apipass=APIPASS    API Passwor for user account.
      --version            Show application version.

➜  jblastor git:(master) ✗
```
