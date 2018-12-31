# README

**SUMMARY**: JBlastor is a command line utility written in go for parsing JSON
files then performing a POST request to the target REST API endpoint.


## USAGE

```
➜  jblastor git:(master) ✗ ./jblastor --help
usage: jblastor --files=FILES [<flags>] <endpoint> [<random>]

Flags:
      --help         Show context-sensitive help (also try --help-long and
                     --help-man).
      --debug        Enable debug mode.
  -t, --timeout=5s   Timeout waiting for POST request.
      --files=FILES  Path to file or directory of file(s) to parse and POST.
      --randomize    Enable randomization of data in JSON files.
      --version      Show application version.

Args:
  <endpoint>  REST API endpoint to send request to.
  [<random>]  Number of random POST requests, per found file, to send [defaults to
              name].

➜  jblastor git:(master) ✗ 

```
