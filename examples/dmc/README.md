# dmc - Sample program on getting a Delegated Machine Credential (DMC) token

## Build sample program

1. git clone https://github.com/centrify/platform-sdk
2. cd platform-sdk/examples/dmc
3. go build ./...

The resulting executable can be found in the current directory.  You can also run "go install ./..." afterwards which installs the executable to the target bin directory in your go environment. 
## Centrify Client requirements

1. Install Centrify Client (version 21.5 or later) on the system.
2. Enroll Centrify Client to a PAS tenant:
  - enable DMC feature by specifying "-F all" or "-F dmc" in cenroll command line.
  - specify the DMC scope(s) by specifying -d in cenroll command line.
  
### Example:

cenroll -t ${MY_TENANT_URL} -c ${MY_ENROLLEMNT_CODE} -F dmc -d testsdk:security/whoami

Notes:

1. MY_TENANT_URL is an environment variable that has the tenant URL
2. MY_ENROLLMENT_CODE is an environment variable that has the enrollment code.

## Usage

1. You need to run this program as root in Linux or a privileged user in Windows.
2. You need to specify a DMC scope using -scope parameter.

### Example:

sudo dmc -scope testsdk

## Sample Linux shell script

This is a simple Linux shell script that demonstrates how to get a DMC token and use it to send a REST API to PAS.

    #!/bin/bash
    dmc=/usr/local/bin/dmc
    OTOKEN=$(sudo $dmc -scope testsdk)
    if [ $? -ne 0 ]
    then
        echo "Cannot get DMC token"
    exit $?
    fi
    TENANT=$(cinfo -T)
    if [ $? -ne 0 ] 
    then 
        echo "Cannot get tenant information"
        exit $?
    fi
    curl -H "Authorization:Bearer $OTOKEN" $TENANT/security/whoami

You can customize this script for:
- different locations of this sample program.  This script assumes that it is already copied to /usr/local/bin.
- the DMC scope name
- the REST API and payload
