package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug     = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout   = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("5s").OverrideDefaultFromEnvar("PING_TIMEOUT").Short('t').Duration()
	files     = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Required().String()
	randomize = kingpin.Flag("randomize", "Enable randomization of data in JSON files.").Bool()
	endpoint  = kingpin.Arg("endpoint", "REST API endpoint to send request to.").Required().String()
	keys      = kingpin.Arg("random", "Number of random POST requests, per found file, to send [defaults to name].").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	fmt.Printf("Would parse file(s): %s to endpoint %s, with timeout %s \n", *files, *endpoint, *timeout)
	fmt.Printf("Randomize is set to: %s \n", *randomize)
	fmt.Printf("Following keys will have randomized values: %s \n", *keys)
}
