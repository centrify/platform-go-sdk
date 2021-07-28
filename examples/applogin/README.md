# applogin - Sample program on getting a HashiCorp Vault token

## Build sample program

1. git clone https://github.com/centrify/platform-go-sdk
2. cd platform-go-sdk/examples/applogin
3. go build ./applogin.go

The resulting executable can be found in the current directory.  You can also run "go install ./..." afterwards which installs the executable to the target bin directory in your go environment. 
## Centrify Client requirements

1. Install Centrify Client (version 21.6 or later) on the system.
2. Enroll Centrify Client to a PAS tenant:
  - enable DMC feature by specifying "-F all" or "-F dmc" in cenroll command line.
  - specify the DMC scope(s) by specifying -d in cenroll command line.
  
### Example:
```
cenroll -t ${MY_TENANT_URL} -c ${MY_ENROLLEMNT_CODE} -F dmc \ 
    -d testsdk:security/whoami \ 
    -d testsdk:usermgmt/getusersrolesandadministrativerights \
    -d 'testsdk:secrets/.*' 
    -d 'testsdk:privilegeddata/.*' 
```
Notes:

1. MY_TENANT_URL is an environment variable that has the tenant URL
2. MY_ENROLLMENT_CODE is an environment variable that has the enrollment code.

## Usage

1. You need to run this program as root in Linux or a privileged user in Windows.
2. You need to specify a DMC scope using -scope parameter.
3. You need to specify a HashiCorp Vault url using -url parameter.

### Example:

```
sudo applogin -scope testsdk -url=http://localhost:8200
```

## Sample Linux shell script

This is a simple Linux shell script that demonstrates how to get a HashiCorp Vault token and use it to send a REST API call to the Vault.

    #!/bin/bash
    program=/usr/local/bin/applogin

    OTOKEN=$(sudo $program -scope testsdk -url=http://localhost:8200)
    if [ $? -ne 0 ]
    then
        echo "Cannot get HashiCorp Vault token"
    exit $?
    fi

    curl -H "X-Vault-Token: $OTOKEN" -X GET http://localhost:8200/v1/centrify/secret_name
    VAULT_TOKEN="$OTOKEN" vault write centrify/secrets/secret_name vault=secret_value


You can customize this script for:
- different locations of this sample program.  This script assumes that it is already copied to /usr/local/bin.
- the DMC scope name
- the HashiCorp Vault REST API URL and payload
