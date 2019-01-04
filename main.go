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
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug     = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout   = kingpin.Flag("timeout", "Timeout waiting for POST request.").Default("15s").Short('t').Duration()
	files     = kingpin.Flag("files", "Path to file or directory of file(s) to parse and POST.").Short('f').Required().ExistingFileOrDir()
	randomize = kingpin.Flag("randomize", "Enable randomization of data in JSON files.").Short('r').Bool()
	endpoint  = kingpin.Flag("endpoint", "REST API endpoint to send request to.").Short('e').Required().String()
	keys      = kingpin.Flag("keys", "Provide list of keys to be randomized [defaults to name].").PlaceHolder("HOSTNAME").Default("name").String()
	randCount = kingpin.Flag("randCount", "Number of random POST request, per found file, to send.").Short('c').Default("1").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Example: jblastor --files /usr/local/myfile.json --endpoint 'http://localhost:8088/save' --randomize --keys 'hostname,name,mc_cfg'"
	kingpin.Parse()
	parsedKeys := parseKeys(*keys)
	processedFiles := processFiles(*files)

	// May remove these once all things are settled
	if *debug {
		fmt.Printf("Debug: will parse file(s): %v \n", *files)
		fmt.Printf("Debug: will perform POST request to: %v with a timeout of: %v \n", *endpoint, *timeout)
		fmt.Printf("Debug: randomize is set to: %v \n", *randomize)
		fmt.Printf("Debug: randomize Count is: %v \n", *randCount)
		fmt.Printf("Debug: following keys will have randomized values: %s \n", parsedKeys)
		fmt.Println("Debug: following files will be processed: ", processedFiles)
	}

	// TEMP: to help with figuring things out.
	prepJSONFile(processedFiles)

	// Next Steps:
	//  - pass processedFiles to function
	//  - that function will take each file, randomize key values if needed and
	//  - concurrenty perform a POST to our target *endpoint
	//
	// First, focus on doing a POST of each file, vanilla, to the endpoint.
	// Once that is working properly, then look to add randomization and finally
	// add concurrency.
}

func postJSON(jsonData []byte) {
	fmt.Printf("DEBUG: url we will post JSON to: %v with a timeout of: %v \n", *endpoint, *timeout)

	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

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
	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(byteValue))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("perfapi", "f6cd3459f9a39c9784b3e328f05be0f7")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	// --------------------------------------------------------------

	// TRIAL - substitute a key/value **NOTE: THIS WORKS!!**
	//result["hostname"] = "somerandom-host"

	//fmt.Println("Processed JSON for host: ", result["hostname"])
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

// TODO: need to build up/off of this.
func prepJSONFile(xf []string) {
	for _, file := range xf {
		readJSONFile(file)
	}
}
