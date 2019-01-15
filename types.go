package sc

import (
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

/*
This file contains various utility types and other internal
structures used throughout. SecurityCenter data types,
(i.e. scans, assets, users, groups, etc.) can be found
defined in their respective files.
*/

// SC is the request multiplexor for various SC API requsts
type SC struct {
	Host string
	Keys map[string]string
}

// Request is the request factory for raw requests to SC
type Request struct {
	sc     *SC
	method string
	path   string
	data   map[string]interface{}
}

// scResp is the type used to marshal API responses from SecurityCenter
type scResp struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"type"`
	Warnings  interface{} `json:"warnings,omitempty"`
}

// Result is the result returned from a raw request to SC
type Response struct {
	// Status code of the HTTP request made
	Status int
	// URL of the HTTP request made
	URL string
	// Data response from the API in simplejson/json format
	Data *simplejson.Json
	// RespBody is the []byte of the API response body
	RespBody []byte
	// HTTPResp is the raw net/http request for easy request customizations
	HTTPResp *http.Response
}
