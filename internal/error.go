package internal

import (
	"encoding/json"
)

type Error struct {
	Status  string `json:"error_status"`
	Message string `json:"error_details"`
	Code    int    `json:"error"`
}

func DecodeError(data []byte) *Error {
	var e Error

	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil
	}

	return &e
}

func (e *Error) String() string {
	return e.Status + ": " + e.Message
}
