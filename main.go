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
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug    = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout  = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("15s").Short('t').Duration()
	files    = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Short('f').Required().ExistingFileOrDir()
	endpoint = kingpin.Flag("endpoint", "REST API endpoint to send request to.").Short('e').Required().String()
	apiUser  = "perfapi"
	apiPass  = "f6cd3459f9a39c9784b3e328f05be0f7"

	//apiUser = kingpin.Flag("apiuser, "API User account permitted to do POST requests.").Short('u').String()
	//apiPass = kingpin.Flag("apipass, "API Passwor for user account.").Short('p').String()
)

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

	// TEMP: to help with figuring things out.
	postJSONFiles(processedFiles)

	// ch := make(chan string)
	// for _,file := range processedFiles{
	// 	go MakeRequest(file, ch)
	// }

	// Next Steps:
	//  - pass processedFiles to function
	//  - concurrenty perform a POST to our target *endpoint
	//
	// First, focus on doing a POST of each file, vanilla, to the endpoint.
	// Once that is working properly, then look to add concurrency.
}

// func postJSON(jsonData []byte) {
// 	fmt.Printf("DEBUG: url we will post JSON to: %v with a timeout of: %v \n", *endpoint, *timeout)

// 	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(jsonData))
// 	req.Header.Set("X-Custom-Header", "myvalue")
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("response Status:", resp.Status)
// 	fmt.Println("response Headers:", resp.Header)
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Println("response Body:", string(body))
// }

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

func readJSONFile(file string) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	// --------------------------------------------------------------
	// TRIAL::
	// --------------------------------------------------------------
	// ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), url)

	start := time.Now()
	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(byteValue))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apiUser, apiPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	secs := time.Since(start).Seconds()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Printf("%.2f seconds elapsed", secs)
	// --------------------------------------------------------------

	// Pretty print the JSON
	if *debug {
		ppJSON, _ := json.MarshalIndent(result, "", "\t")
		if ppJSON != nil {
			fmt.Println(string(ppJSON))
		}
	}
}

// TODO: need to build up/off of this.
func postJSONFiles(xf []string) {
	for _, file := range xf {
		readJSONFile(file)
	}
}
