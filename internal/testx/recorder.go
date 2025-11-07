package testx

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

type Response struct {
	*httptest.ResponseRecorder
}

func NewRecorder() *Response           { return &Response{ResponseRecorder: httptest.NewRecorder()} }
func (r *Response) BodyString() string { return r.Body.String() }

func NewRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
