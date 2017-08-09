package api

import (
	"time"
)

type Response struct {
	Success   bool        `json:"success"`
	Errors    []Error    `json:"errors,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp int64       `json:"ts,omitempty"`
}

var EmptyPayload interface{}

func SuccessResponse(payload interface{}) *Response {
	return &Response{
		Success:   true,
		Payload:   payload,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
}

func ErrorResponse(errors []Error) *Response {
	return &Response{
		Success:   false,
		Errors:    errors,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
}
