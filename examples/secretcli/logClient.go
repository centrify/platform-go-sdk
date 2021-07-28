package main

import (
	"net/http"
)

// newLogClient returns a new http client object but using a different transport
// that logs messages
func newLogClient() *http.Client {
	result := &http.Client{} // DO NOT use DefaultClient here since we don't want to modify it

	// create new transport
	origTransport := result.Transport
	if origTransport == nil {
		origTransport = http.DefaultTransport
	}
	newTransport := &logTransport{
		xprt: origTransport,
	}
	result.Transport = newTransport
	return result
}
