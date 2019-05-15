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

	logr "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
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

var ch = make(chan HTTPResponse)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:8088/save' "
	kingpin.Parse()
	processedFiles := processFiles(*files)

	// playing with LogRus
	logr.SetFormatter(&logr.JSONFormatter{})
	logr.SetOutput(os.Stdout)
	logr.SetLevel(logr.InfoLevel)

	fileCount := len(processedFiles)
	logr.Info("Number of files: ", fileCount)

	for _, file := range processedFiles {
		// For each URL call the DOHTTPPost function (concurrency)
		// TODO: need to add some throttling.  Seems the speed of requests going to
		//       the target endpoint is causing a lock up. (05/09/19)
		//go DoHTTPPost(file, ch)
		DoHTTPPost(file)
	}

	// for range processedFiles {
	// 	// Want to use both of these but in a log file.
	// 	fmt.Println((<-ch).status) // example:  201 Created
	// 	//fmt.Println(string((<-ch).body))
	// }
}

// DoHTTPPost function takes a path to a JSON formatted file, extracts the JSON
// data and then does a POST request to the targetted endpoint, concurrently.
//func DoHTTPPost(file string, ch chan<- HTTPResponse) {
func DoHTTPPost(file string) {
	jsonFile, err := os.Open(file)
	if err != nil {
		// Change this to log error to file.
		logr.Warn("DoHTTPPost: Error opening file")
		return
	}
	defer jsonFile.Close()

	jsonValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		// Change this to log error to file.
		logr.Warn("DoHTTPPost: Error reading JSON file.")
		return
	}

	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(jsonValue))
	req.Header.Set("X-Custom-Header", "JBLASTOR")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(*apiUser, *apiPass)

	client := &http.Client{}
	httpResponse, err := client.Do(req)

	if err != nil {
		// Change this to log error to file or ignore
		logr.Warn("DoHTTPPost#httpResponse: error making http request")
		return
	}
	defer httpResponse.Body.Close()

	httpBody, _ := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		// Change this to log error to file or ignore
		logr.Warn("DoHTTPPost#httpBody: error ready response body")
		return
	}
	r := HTTPResponse{httpResponse.Status, httpBody}
	fmt.Printf("Status: %v | Data: %v \n", r.status, string(r.body))
	// Look into converting above to use Logrus (log with fields)
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
