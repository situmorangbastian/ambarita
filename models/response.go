package models

import (
	"encoding/json"
	"net/http"
)

// Response represents http response
type Response struct {
	Status   int                     `json:"status"`
	Message  string                  `json:"message"`
	Data     interface{}             `json:"data"`
	PageInfo *map[string]interface{} `json:"page_info,omitempty"`
}

// MarshalJSON converts response to JSON bytes.
func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Status   int                     `json:"status"`
		Message  string                  `json:"message"`
		Data     interface{}             `json:"data"`
		PageInfo *map[string]interface{} `json:"page_info,omitempty"`
	}{
		Status:   r.Status,
		Message:  r.Message,
		Data:     r.Data,
		PageInfo: r.PageInfo,
	})
}

// DefaultErrorResponse return default value for Response
func DefaultErrorResponse() Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}
}

// DefaultSuccessResponse return default value for Response
func DefaultSuccessResponse() Response {
	return Response{
		Status:  http.StatusOK,
		Message: "Success",
	}
}

// DefaultCreatedResponse return default value for Response
func DefaultCreatedResponse() Response {
	return Response{
		Status:  http.StatusCreated,
		Message: "Created",
	}
}
