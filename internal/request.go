package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

// StreamPrint opens a stream and print it in a goroutine
func StreamPrint(method string, path string, args []byte, mods ...RequestModifier) {
	reader, _, err := Stream(method, path, args, mods...)

	if err != nil {
		Exit("Error while attach: %s\n", err)
	}

	go func(stream io.ReadCloser) {
		DisplayStream(reader)
	}(reader)
}

// ReqWant requests with a method on a path, check wantCode and returns []byte
func ReqWant(method string, wantCode int, path string, jsonStr []byte) []byte {
	return apiRequest(method, wantCode, path, jsonStr, false)
}

// apiRequest helper, issue the request and consume stream if requested. Otherwise, return full body
func apiRequest(method string, wantCode int, path string, jsonStr []byte, stream bool) []byte {
	bodyStream, code, err := doRequest(method, path, jsonStr)
	Check(err)

	defer bodyStream.Close()

	if stream && code == wantCode {
		DisplayStream(bodyStream)
		return nil
	}

	var body []byte
	body, err = ioutil.ReadAll(bodyStream)
	Check(err)

	if code != wantCode {
		if err == nil {
			FormatOutputDef(body)
		}
		os.Exit(1)
	}

	if Verbose {
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(body))
	}
	return body

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

	respBody, code, err := doRequest(method, path, args, mods...)
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
	return doRequest(method, path, args, mods...)
}

//doRequest builds the request and return io.ReadCloser
func doRequest(method string, path string, args []byte, mods ...RequestModifier) (io.ReadCloser, int, error) {

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

	for key, val := range Headers {
		if Verbose {
			fmt.Fprintf(os.Stderr, "Request header: %s=%s\n", key, val)
		}
		req.Header.Set(key, val)
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

// DisplayStream decode each line from http buffer and print either message or error. Return last read line
func DisplayStream(buffer io.ReadCloser) ([]byte, error) {
	reader := bufio.NewReader(buffer)

	for {
		line, err := reader.ReadBytes('\n')
		if Verbose {
			fmt.Fprintf(os.Stderr, "Received message: %v\n", string(line))
		}

		// Progress message
		m := DecodeMessage(line)
		if m != nil {
			fmt.Fprintln(os.Stderr, m.Message)
			continue
		}

		// Error message (will be last message)
		e := DecodeError(line)
		if e != nil {
			return line, fmt.Errorf(e.Message)
		}

		// Final message
		if err == io.EOF {
			return line, nil
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
func Check(err error) {
	if err != nil {
		if Verbose {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
