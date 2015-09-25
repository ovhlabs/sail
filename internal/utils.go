package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
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
