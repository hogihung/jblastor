package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug    = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout  = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("15s").Short('t').Duration()
	files    = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Short('f').Required().ExistingFileOrDir()
	endpoint = kingpin.Flag("endpoint", "REST API endpoint to send request to.").Short('e').Required().String()
	apiUser  = kingpin.Flag("apiuser", "API User account permitted to do POST requests.").Short('u').String()
	apiPass  = kingpin.Flag("apipass", "API Passwor for user account.").Short('p').String()
)

// HTTPResponse is a struct for handling the responses we will be getting from
// the POST requests.
type HTTPResponse struct {
	status string
	body   []byte
}

var ch chan HTTPResponse = make(chan HTTPResponse)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:8088/save' "
	kingpin.Parse()
	processedFiles := processFiles(*files)

	// May remove these once all things are settled
	if *debug {
		fmt.Printf("Debug: will parse file(s): %v \n", *files)
		fmt.Printf("Debug: will perform POST request to: %v with a timeout of: %v \n", *endpoint, *timeout)
		fmt.Println("Debug: following files will be processed: ", processedFiles)
	}

	for _, file := range processedFiles {
		//For each URL call the DOHTTPPost function (notice the go keyword)
		go DoHTTPPost(file, ch)
	}

	for range processedFiles {
		// Use the response (<-ch).body
		fmt.Println((<-ch).status)
	}
}

// DoHTTPPost function takes a path to a JSON formatted file, extracts the JSON
// data and then does a POST request to the targetted endpoint, concurrently.
func DoHTTPPost(file string, ch chan<- HTTPResponse) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonValue, _ := ioutil.ReadAll(jsonFile)

	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(jsonValue))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apiUser, apiPass)

	client := &http.Client{}
	httpResponse, err := client.Do(req)
	if err != nil {
		panic(err)
		//fmt.Printf("Error encountered with file %v.  Error is: %v. \n", file, err)
		//fmt.Println("Error with file: ", file)
	}
	defer httpResponse.Body.Close()

	httpBody, _ := ioutil.ReadAll(httpResponse.Body)
	ch <- HTTPResponse{httpResponse.Status, httpBody}
}

func processFiles(f string) []string {
	parsedFiles := make([]string, 0)

	directory, err := isDirectory(f)
	if err != nil {
		fmt.Println(err)
	}

	if directory {
		files, err := ioutil.ReadDir(f)
		if err != nil {
			log.Fatal(err)
		}

		var filename string
		for _, file := range files {
			filename = f + file.Name()
			jsonFile := isJSONFile(filename)
			if jsonFile == false || err != nil {
				continue
			}
			parsedFiles = append(parsedFiles, filename)
		}
	} else {
		jsonFile := isJSONFile(f)
		if jsonFile == false || err != nil {
			fmt.Printf("File %v is not a valid JSON file. \n", f)
		} else {
			parsedFiles = append(parsedFiles, f)
		}

	}
	return parsedFiles
}

func isJSONFile(file string) bool {
	if (filepath.Ext(file) == ".json") && isValidJSON(file) {
		return true
	}
	return false
}

func isValidJSON(file string) bool {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if json.Valid(byteValue) {
		return true
	}
	return false
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
