package main

import (
	"fmt"
	"net/http"
	"time"
)

// logTransport is a custom RoundTripper object that supports the
// http.RoundTripper interface.
// Its implementation of RoundTrip() logs the time that the REST API
// is sent and time received, as well as the REST call transaction ID (xid)
type logTransport struct { // whether to dump full HTTP request and response
	xprt   http.RoundTripper // use this roundtripper to send/receive the original REST API
	xid    string
	status int
	err    error
}

// start message: current_time|starts|METHOD|URL
const logStartMsg = "%v REST|starts|%s|%s\n"

// end message: current_time|ends|METHOD|URL|HTTP status code|xid|elapsed time
const logEndMsg = "%v REST|ends|%s|%s|%d|%s|%v\n"

// error message: current_time|error|METHOD|URL|elapsed time|error
const logErrMsg = "%v REST|error|%s|%s|%v|%v\n"

func (c *logTransport) logStart(req *http.Request) {
	fmt.Printf(logStartMsg, time.Now().UTC().Format("2006/01/02 15:04:05.000000"),
		req.Method, req.URL.Path)
}

func (c *logTransport) logEnd(req *http.Request, starttime time.Time) {
	if c.err == nil {
		fmt.Printf(logEndMsg, time.Now().UTC().Format("2006/01/02 15:04:05.000000"),
			req.Method, req.URL.Path, c.status, c.xid, time.Since(starttime))
	} else {
		fmt.Printf(logErrMsg, time.Now().UTC().Format("2006/01/02 15:04:05.000000"),
			req.Method, req.URL.Path, time.Since(starttime), c.err)
	}
}

// RoundTrip logs starts/ends messages for the REST API call
func (c *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	starttime := time.Now()
	c.logStart(req)
	defer c.logEnd(req, starttime)

	// get the original roundtripper to do the round trip
	resp, err := c.xprt.RoundTrip(req)
	if err != nil {
		c.err = err
		return resp, err
	}

	// save the transaction ID and status in c so logEnd can pick it up
	c.xid = resp.Header.Get("X-CFY-TX-ID")
	c.status = resp.StatusCode
	return resp, err
}
