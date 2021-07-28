# dmc

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/.)

Package dmc provides application with APIs for obtaining Delegated Machine Credentials (DMC)

## Run go tests

You need to do the followings to have a successful run of go unit tests:

```go
1. Install Centrify Client (version 21.5 or later) on the system.
2. Enroll Centrify Client to a PAS tenant:
  - enable DMC feature by specifying "-F all" or "-F dmc" in cenroll command line.
  - specify the DMC scope "testsdk" by specifying "-d testsdk:security/whoami" in cenroll command line.
3. Run the unit test as root
```

## Sample Program

A sample program can be found in [https://github.com/centrify/platform-go-sdk/examples/dmc](https://github.com/centrify/platform-go-sdk/examples/dmc)

## Functions

### func [GetDMCToken](/dmc.go#L39)

`func GetDMCToken(scope string) (string, error)`

GetDMCToken returns an oauth token for the requested scope that has the
identity of the current machine account.

Possible error returns:

```go
ErrCannotGetToken  - other errors in getting the token
ErrCannotSetupConnection - Cannot setup connection to Centrify Client
ErrClientNotInstalled - Centrify Client is not installed in system
ErrCommunicationError - Communication error with Centrify Client
```

### func [GetEnrollmentInfo](/dmc.go#L103)

`func GetEnrollmentInfo() (string, string, error)`

GetEnrollmentInfo returns information about Centrify Client enrollment information

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
