package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//"os"

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
	fmt.Printf("Would parse file(s): %v to endpoint %s, with timeout %s \n", *files, *endpoint, *timeout)
	fmt.Printf("Randomize is set to: %v \n", *randomize)
	fmt.Printf("Randomize Count is: %v \n", *randCount)
	fmt.Printf("Following keys will have randomized values: %s \n", *keys)

	parsedKeys := parseKeys(*keys)
	parseFiles(*files)
	slurpDir(*files)

	fmt.Printf("Will randomize the values for the keys %v if they exist.", parsedKeys)
}

// Take value from 'files' to build a list of files to be parsed
//
// If randomize is true, validate we have 'keys' -> parse the provided string
// of keys to see which keys we need to randomize for each file provided
//
//

func parseKeys(k string) []string {
	fmt.Println("Debug: the following was passed in to parseKeys: ", k)
	xs := strings.Split(k, ",")
	fmt.Println("Size of xs: ", len(xs))

	// Iterate over each supplied key and randomize the value for supplied key
	// in the target file.
	parsedKeys := make([]string, len(xs))
	for _, value := range xs {
		//fmt.Println("Original value: ", value)
		//fmt.Println("Downcase version of value: ", strings.ToLower(value))
		parsedKeys = append(parsedKeys, strings.ToLower(value))
	}
	fmt.Println("Parsed Keys:", parsedKeys)
	return parsedKeys
}

func parseFiles(f string) {
	fmt.Println("Debug: the following was passed in to parseFiles: ", f)

	// if f is a directory, iterate through all JSON (*.json) files
	// for _, value := range f {
	//  	processFile(f)
	// }

	//
	// if f is a file, process that file
	processFile(f)
}

func processFile(f string) {
	fmt.Println("Debug: processing file: ", f)

	isDir, err := IsDirectory(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Is this a directory? (IsDirectory) ", isDir)
	fmt.Println("FileExists?", FileExists(f))
	fmt.Println("DirExists?", DirExists(f))
}

// IsDirectory comment here
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// NOTES:  Using os.Stat(f) we can leverage the functions IsDir to see if f is
//         a directory, OR use IsRegular to see if f is a file.
//           - if fi.Mode().IsDir()
//           - if fi.Mode().IsRegular()

// Usefull blog:  https://flaviocopes.com/go-list-files/

// FileExists reports whether the named file exists as a boolean
func FileExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// DirExists reports whether the dir exists as a boolean
func DirExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

// trying something - works pretty good!
func slurpDir(d string) {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
