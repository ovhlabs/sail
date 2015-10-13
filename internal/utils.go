package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
)

// Exit func display an error message on stderr and exit 1
func Exit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

// Message type
type Message struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// DecodeMessage return a Message struct from json
func DecodeMessage(data []byte) *Message {
	var m Message

	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil
	}

	if m.Message == "" && m.Type == "" {
		return nil
	}
	return &m
}

// ParseResourceName normalizes repo or service name of the form [[cluster/]application/]name[:tag]
func ParseResourceName(repositoryName string) (host, application, repository, tag string, err error) {
	// FIXME: duplicate run
	ReadConfig()
	host = ""
	application = User
	repository = ""
	tag = ""
	err = nil

	// Split namespace and repository
	split := strings.Split(repositoryName, "/")
	if len(split) == 1 {
		repository = split[0]
	} else if len(split) == 2 {
		application = split[0]
		repository = split[1]
	} else if len(split) == 3 {
		host = split[0]
		application = split[1]
		repository = split[2]
	} else {
		err = fmt.Errorf("Invalid repository %s. Should be of form [[endpoint/]application/]name[:tag]", repositoryName)
		return
	}

	// Split repo URL and tag
	split = strings.Split(repository, ":")
	if len(split) == 2 {
		repository = split[0]
		tag = split[1]
	} else if len(split) > 2 {
		err = fmt.Errorf("Invalid repository %s. Should be of form [[endpoint/]application/]name[:tag]", repositoryName)
		return
	}

	return
}

// CheckName validates that a service or application looks consistent. It will only block request that will *always* fail. It does not duplicate API validation
func CheckName(name string) error {
	if strings.Contains(name, "/") {
		return fmt.Errorf("Name %s can not contain '/'", name)
	}
	return nil
}

// CheckHostConsistent with config. Assume config has already been parsed. Consider '' as OK
func CheckHostConsistent(host string) bool {
	if host == "" {
		return true
	}

	url, err := url.ParseRequestURI(Host)
	if err != nil {
		return false
	}

	return host == url.Host
}

// ExitAfterCtrlC will exit(0) as soon as Ctrl-C is pressed. Typically used when streaming console
func ExitAfterCtrlC() {
	var endWaiter sync.WaitGroup
	var signalChannel chan os.Signal

	endWaiter.Add(1)

	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}
