package http

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"time"
)

// Header represents the HTTP header
type Header struct {
	Key   string
	Value string
}

type RequestOptions struct {
	// URI specifies the request's URI
	URI string

	// Method is the HTTP method
	Method string

	// ContentType is the HTTP content type
	ContentType string

	// JSONBody is the JSON body of the request
	JSONBody interface{}

	// Headers is the HTTP headers
	Headers []Header

	// Timeout is the HTTP timeout in seconds
	Timeout time.Duration
}

// Get request service through http Get method
func Get(uri string, headers ...Header) ([]byte, error) {
	return Request(&RequestOptions{
		URI:     uri,
		Method:  fasthttp.MethodGet,
		Headers: headers,
	})
}

// PostJSON post json formatted request to service via http POST method
func PostJSON(uri string, body interface{}, headers ...Header) ([]byte, error) {
	return Request(&RequestOptions{
		URI:         uri,
		Method:      fasthttp.MethodPost,
		Headers:     headers,
		ContentType: "application/json; charset=utf-8",
		JSONBody:    body,
	})
}

// Request represents a request to a service endpoint
func Request(opts *RequestOptions) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	c := &fasthttp.Client{}
	var err error = nil

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(opts.URI)
	req.Header.SetMethod(opts.Method)

	if opts.ContentType != "" {
		req.Header.SetContentType(opts.ContentType)
	}

	if opts.JSONBody != nil {
		if body, err := json.Marshal(opts.JSONBody); err != nil {
			return nil, err
		} else {
			req.SetBody(body)
		}
	}

	for _, header := range opts.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	if opts.Timeout > 0 {
		err = c.DoTimeout(req, resp, opts.Timeout)
	} else {
		err = c.Do(req, resp)
	}

	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
