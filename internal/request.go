package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func initRequest(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sailabove sailgo CLI/"+VERSION)
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{}
	return &http.Client{Transport: tr}
}

func GetWantJSON(path string) string {
	return GetJSON(ReqWant("GET", http.StatusOK, path, nil))
}

func PostWantJSON(path string) string {
	return GetJSON(ReqWant("POST", http.StatusCreated, path, nil))
}

func DeleteWantJSON(path string) string {
	return GetJSON(ReqWant("DELETE", http.StatusOK, path, nil))
}

func ReqWantJSON(method string, wantCode int, path string, body []byte) string {
	return GetJSON(ReqWant(method, wantCode, path, body))
}

func StreamWant(method string, wantCode int, path string, jsonStr []byte) {
	ApiRequest(method, wantCode, path, jsonStr, true)
}

func ReqWant(method string, wantCode int, path string, jsonStr []byte) []byte {
	return ApiRequest(method, wantCode, path, jsonStr, false)
}

func ApiRequest(method string, wantCode int, path string, jsonStr []byte, stream bool) []byte {

	err := ReadConfig()
	if err != nil {
		fmt.Printf("Error reading configuration: %s\n", err)
		os.Exit(1)
	}

	var req *http.Request
	if jsonStr != nil {
		req, _ = http.NewRequest(method, Host+path, bytes.NewReader(jsonStr))
	} else {
		req, _ = http.NewRequest(method, Host+path, nil)
	}

	initRequest(req)
	req.SetBasicAuth(User, Password)
	resp, err := getHTTPClient().Do(req)
	Check(err)
	defer resp.Body.Close()

	var body []byte
	if !stream {
		body, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode != wantCode || Verbose {
		fmt.Printf("Response Status : %s\n", resp.Status)
		fmt.Printf("Request path : %s\n", Host+path)
		fmt.Printf("Request Headers : %s\n", req.Header)
		fmt.Printf("Request Body : %s\n", string(jsonStr))
		fmt.Printf("Response Headers : %s\n", resp.Header)
		if err == nil {
			fmt.Printf("Response Body : %s\n", string(body))
		}
		if !Verbose {
			os.Exit(1)
		}
	}

	if stream {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return nil
			}
			if string(line) != "" {
				log.Print(string(line))
			}
		}
	} else {
		Check(err)
		return body
	}
}

func GetListApplications(args []string) []string {
	apps := []string{}
	if len(args) == 0 {
		b := ReqWant("GET", http.StatusOK, "/applications", nil)
		err := json.Unmarshal(b, &apps)
		Check(err)
	}
	return apps
}

func GetJSON(s []byte) string {
	if Pretty {
		var out bytes.Buffer
		json.Indent(&out, s, "", "  ")
		return out.String()
	}
	return string(s)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
