package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug     = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout   = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("5s").Short('t').Duration()
	files     = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Short('f').Required().ExistingFileOrDir()
	randomize = kingpin.Flag("randomize", "Enable randomization of data in JSON files.").Short('r').Bool()
	endpoint  = kingpin.Flag("endpoint", "REST API endpoint to send request to.").Short('e').Required().String()
	keys      = kingpin.Flag("keys", "Provide list of keys to be randomized [defaults to name].").PlaceHolder("HOSTNAME").Default("name").String()
	randCount = kingpin.Flag("randCount", "Number of random POST request, per found file, to send.").Short('c').Default("1").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:7888/save' --randomize --keys 'hostname,name,mc_cfg'"
	kingpin.Parse()
	parsedKeys := parseKeys(*keys)
	processedFiles := processFiles(*files)

	prepJSONFile(processedFiles)

	if *debug {
		fmt.Printf("Debug: will parse file(s): %v \n", *files)
		fmt.Printf("Debug: will perform POST request to: %v with a timeout of: %v \n", *endpoint, *timeout)
		fmt.Printf("Debug: randomize is set to: %v \n", *randomize)
		fmt.Printf("Debug: randomize Count is: %v \n", *randCount)
		fmt.Printf("Debug: following keys will have randomized values: %s \n", parsedKeys)
		fmt.Println("Debug: following files will be processed: ", processedFiles)
	}
}

// NOTES:  Need to capture the directory/full path name when slurping the files
//         in a directory, not just the filenames.

// If randomize is true, validate we have 'keys' -> parse the provided string
// of keys to see which keys we need to randomize for each file provided
//
// At this time we will only support randomizing the value of top level keys.
// For example, {"hostname:" "somefqdn-here"} OR {"podname:" "some-pod-name-here"}
func parseKeys(k string) []string {
	keyString := strings.Split(k, ",")

	parsedKeys := make([]string, 0)
	for _, value := range keyString {
		parsedKeys = append(parsedKeys, strings.ToLower(value))
	}
	return parsedKeys
}

func processFiles(f string) []string {
	parsedFiles := make([]string, 0)

	directory, err := IsDirectory(f)
	if err != nil {
		fmt.Println(err)
	}

	if directory {
		fmt.Println("Passed argument is a directory: ", f)

		files, err := ioutil.ReadDir(f)
		if err != nil {
			log.Fatal(err)
		}

		var filename string
		for _, file := range files {
			filename = file.Name()

			jsonFile := isJSONFile(filename)
			if jsonFile == false || err != nil {
				continue
			}

			filename = f + filename
			parsedFiles = append(parsedFiles, filename)
		}
		fmt.Println("Parsed Files: ", parsedFiles)

	} else {
		parsedFiles = append(parsedFiles, f)
	}
	return parsedFiles
}

// Returns path/filename if file has json extension
func isJSONFile(file string) bool {
	if filepath.Ext(file) == ".json" {
		return true
	}
	return false
}

// IsDirectory checks passed in argument to see if it is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func readJSONFile(file string) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	// TRIAL - substitute a key/value **NOTE: THIS WORKS!!**
	//result["hostname"] = "somerandom-host"

	fmt.Println("Processed JSON for host: ", result["hostname"])
	//fmt.Println("Processed JSON asset.properties: ", result["asset.properties"])
	//fmt.Println("Processed JSON mc_cfg: ", result["mc_cfg"])

	// Pretty print the JSON
	if *debug {
		newResult, _ := json.MarshalIndent(result, "", "\t")
		if newResult != nil {
			fmt.Println(string(newResult))
		}
	}
}

func prepJSONFile(xf []string) {
	for _, file := range xf {
		readJSONFile(file)
	}
}
