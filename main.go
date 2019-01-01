package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	if *debug {
		fmt.Printf("Debug: will parse file(s): %v \n", *files)
		fmt.Printf("Debug: will perform POST request to: %v with a timeout of: %v \n", *endpoint, *timeout)
		fmt.Printf("Debug: randomize is set to: %v \n", *randomize)
		fmt.Printf("Debug: randomize Count is: %v \n", *randCount)
		fmt.Printf("Debug: following keys will have randomized values: %s \n", parsedKeys)
		fmt.Println("Debug: following files were passed in: ", processedFiles)
	}
}

// If randomize is true, validate we have 'keys' -> parse the provided string
// of keys to see which keys we need to randomize for each file provided
func parseKeys(k string) []string {
	xs := strings.Split(k, ",")

	parsedKeys := make([]string, 0)
	for _, value := range xs {
		parsedKeys = append(parsedKeys, strings.ToLower(value))
	}
	return parsedKeys
}

func processFiles(f string) []string {
	parsedFiles := make([]string, 0)

	// TODO: Need to only gather files with extension of json (*.json)
	if *debug {
		fmt.Println("Debug: the following was passed in to parseFiles: ", f)
		fmt.Printf("Debug: the type of passed in argument is: %T \n", f)
	}

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

		for _, file := range files {
			name := file.Name()
			parsedFiles = append(parsedFiles, name)

		}
		fmt.Println("Parsed Files: ", parsedFiles)

	} else {
		//fmt.Println("Passed argument is a file: ", f)
		parsedFiles = append(parsedFiles, f)
		//fmt.Println("Parsed File: ", parsedFiles)
		//fmt.Printf("parsedFile is type: %T \n", f)
	}
	return parsedFiles
}

// func processFile(f string) {
// 	if *debug {
// 		fmt.Println("Debug: processing file: ", f)
// 		fmt.Println("Debug: FileExists?", FileExists(f))
// 		fmt.Println("Debug: DirExists?", DirExists(f))
// 	}

// 	isDir, err := IsDirectory(f)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	if *debug {
// 		fmt.Println("Is this a directory? (IsDirectory) ", isDir)
// 	}
// }

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
