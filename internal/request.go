package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func initRequest(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sailabove sail CLI/"+VERSION)
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{}
	return &http.Client{Transport: tr}
}

// GetWantJSON GET on path and return []byte of JSON
func GetWantJSON(path string) []byte {
	return ReqWant("GET", http.StatusOK, path, nil)
}

// PostWantJSON POST on path and return []byte of JSON
func PostWantJSON(path string) []byte {
	return ReqWant("POST", http.StatusCreated, path, nil)
}

// PostBodyWantJSON POST a body on path and return []byte of JSON
func PostBodyWantJSON(path string, body []byte) []byte {
	return ReqWant("POST", http.StatusCreated, path, body)
}

// DeleteWantJSON on path and return []byte of JSON
func DeleteWantJSON(path string) []byte {
	return ReqWant("DELETE", http.StatusOK, path, nil)
}

// DeleteBodyWantJSON on path and return []byte of JSON
func DeleteBodyWantJSON(path string, body []byte) []byte {
	return ReqWant("DELETE", http.StatusOK, path, body)
}

// ReqWantJSON requests with a method on a path, check wantCode and returns []byte of JSON
func ReqWantJSON(method string, wantCode int, path string, body []byte) []byte {
	return ReqWant(method, wantCode, path, body)
}

// StreamWant request a path with method and stream result
func StreamWant(method string, wantCode int, path string, jsonStr []byte) {
	apiRequest(method, wantCode, path, jsonStr, true)
}

// ReqWant requests with a method on a path, check wantCode and returns []byte
func ReqWant(method string, wantCode int, path string, jsonStr []byte) []byte {
	return apiRequest(method, wantCode, path, jsonStr, false)
}

func apiRequest(method string, wantCode int, path string, jsonStr []byte, stream bool) []byte {

	err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration: %s\n", err)
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

	for key, val := range Headers {
		if Verbose {
			fmt.Fprintf(os.Stderr, "Request header: %s=%s\n", key, val)
		}
		req.Header.Set(key, val)
	}

	resp, err := getHTTPClient().Do(req)
	Check(err)
	defer resp.Body.Close()

	var body []byte
	if !stream {
		body, err = ioutil.ReadAll(resp.Body)
	}

	if Verbose {
		fmt.Fprintf(os.Stderr, "Request path: %s\n", Host+path)
		fmt.Fprintf(os.Stderr, "Request Headers: %s\n", req.Header)
		fmt.Fprintf(os.Stderr, "Request Body: %s\n", string(jsonStr))
		fmt.Fprintf(os.Stderr, "Response Headers: %s\n", resp.Header)
		fmt.Fprintf(os.Stderr, "Response Status: %s\n", resp.Status)

		if err == nil {
			fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))

		}
	}

	if resp.StatusCode != wantCode {
		if err == nil {
			FormatOutputDef(body)
		}
		os.Exit(1)
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

// RequestModifier is used to modify behavior of Request and Steam functions
type RequestModifier func(req *http.Request)

// SetHeader modify headers of http.Request
func SetHeader(key, value string) RequestModifier {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// Request executes an authentificated HTTP request on $path given $method and $args
func Request(method string, path string, args []byte, mods ...RequestModifier) ([]byte, int, error) {

	respBody, code, err := Stream(method, path, args, mods...)
	if err != nil {
		return nil, 0, err
	}
	defer respBody.Close()

	var body []byte
	body, err = ioutil.ReadAll(respBody)
	if err != nil {
		return nil, 0, err
	}

	if Verbose {
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", body)
	}

	return body, code, nil
}

// Stream makes an authenticated http request and return io.ReadCloser
func Stream(method string, path string, args []byte, mods ...RequestModifier) (io.ReadCloser, int, error) {

	err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration: %s\n", err)
		os.Exit(1)
	}

	var req *http.Request
	if args != nil {
		req, _ = http.NewRequest(method, Host+path, bytes.NewReader(args))
	} else {
		req, _ = http.NewRequest(method, Host+path, nil)
	}
	initRequest(req)

	for i := range mods {
		mods[i](req)
	}

	req.SetBasicAuth(User, Password)
	resp, err := getHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}

	if Verbose {
		fmt.Fprintf(os.Stderr, "Response Status: %s\n", resp.Status)
		fmt.Fprintf(os.Stderr, "Request path: %s\n", Host+path)
		fmt.Fprintf(os.Stderr, "Request Headers: %s\n", req.Header)
		fmt.Fprintf(os.Stderr, "Request Body: %s\n", string(args))
		fmt.Fprintf(os.Stderr, "Response Headers: %s\n", resp.Header)
	}

	return resp.Body, resp.StatusCode, nil
}

// DisplayStream decode each line from http buffer and print either message or error
func DisplayStream(buffer io.ReadCloser) error {
	reader := bufio.NewReader(buffer)

	for {
		line, err := reader.ReadBytes('\n')
		m := DecodeMessage(line)
		if m != nil {
			fmt.Fprintln(os.Stderr, m.Message)
		}
		e := DecodeError(line)
		if e != nil {
			return e
		}
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// GetListApplications returns list of applications, GET on /applications
func GetListApplications(args []string) []string {
	apps := []string{}
	if len(args) == 0 {
		b := ReqWant("GET", http.StatusOK, "/applications", nil)
		err := json.Unmarshal(b, &apps)
		Check(err)
	}
	return apps
}

// Check checks e and panic if not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
