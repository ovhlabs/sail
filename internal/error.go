package internal

import (
	"encoding/json"
)

// Error type
type Error struct {
	Status  string `json:"error_status"`
	Message string `json:"error_details"`
	Code    int    `json:"error"`
}

// DecodeError return an Error struct from json
func DecodeError(data []byte) *Error {
	var e Error

	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil
	}

	if e.Message == "" && e.Status == "" {
		return nil
	}
	return &e
}

func (e *Error) String() string {
	return e.Status + ": " + e.Message
}

func (e *Error) Error() string {
	return e.String()
}
