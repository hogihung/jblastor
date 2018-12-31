package main

import (
	"fmt"
	//"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug     = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout   = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("5s").Short('t').Duration()
	files     = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Short('f').Required().ExistingFileOrDir()
	randomize = kingpin.Flag("randomize", "Enable randomization of data in JSON files.").Short('r').Bool()
	endpoint  = kingpin.Flag("endpoint", "REST API endpoint to send request to.").Short('e').Required().String()
	keyz      = kingpin.Flag("keys", "Provide list of keys to be randomized [defaults to name].").PlaceHolder("HOSTNAME").Default("name").String()
	randCount = kingpin.Flag("randCount", "Number of random POST request, per found file, to send.").Short('c').Default("1").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:7888/save' --randomize --keys 'hostname,name,mc_cfg'"
	kingpin.Parse()
	fmt.Printf("Would parse file(s): %v to endpoint %s, with timeout %s \n", *files, *endpoint, *timeout)
	fmt.Printf("Randomize is set to: %v \n", *randomize)
	fmt.Printf("Following keys will have randomized values: %s \n", *keyz)

	parseKeys(*keyz)
	parseFiles(*files)
}

// Take value from 'files' to build a list of files to be parsed
//
// If randomize is true, validate we have 'keys' -> parse the provided string
// of keys to see which keys we need to randomize for each file provided
//
//

func parseKeys(k string) {
	fmt.Printf("Debug: the following %v was passed in to parseKeys", k)
}

func parseFiles(f string) {
	fmt.Printf("Debug: the following %v was passed in to parseFiles", f)
}
