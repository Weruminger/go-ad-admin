package kea

import "net/http"

type Client struct {
	HTTP  *http.Client
	URL   string
	Token string
}
