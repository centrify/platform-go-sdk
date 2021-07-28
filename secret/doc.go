/*
Package secret manages secrets stored in Centrify PAS using a simple
set of APIs.

Obtain a client handle

The application must call NewSecretClient() to obtain a client handle to Centrify PAS.

Access credential

You can specify an OAuth access token in NewSecretClient().  It will be used for all subsequent calls
to access the secrets.  Altenatively, you can setup the necessary authorization header by calling
the method AddDefaultHeaders().

Custom HTTP Client

You can specify to use a custom HTTP client by providing a HTTPFactory in NewSecretClient().  If this is not
specified, the default http.DefaultClient is used.  One example of using a custom HTTP client is for additional logging
of REST calls.

See also

The file logClient.go in github.com/centrify/platform-go-sdk/examples/secretcli is an example
of using a custom HTTP client.

Accessing secrets

The following methods are provided:

  Create:       Create a secret
  CreateFolder: Create a secret folder
  Delete:       Delete a secret/folder
  Get:          Get value of a secret
  GetMetaData:  Get metadata associated with a secret
  List:         List secrets in a secret folder
  Modify:       Modify a secert

Additional customizations

  AddDefaultHeaders:    Add additional HTTP header(s) to each outgoing HTTP request.
  SetDebug:             Enable logging of debug messages.  Debug should be OFF (default)
                        in production environment as the secret may be shown in logs.
  SetUserAgent:         Set the UserAgent header.

Tenant setup to run go tests

You need to create a web application in PAS.  You need to setup the following parameters about
the web application:

- Settings
  * Application ID (e.g., sdktest)
  * Name (e.g., sdktest)

- General Usage
  * Client ID Type: Confidential
  * "Must be OAuth Client" checkbox should be left unchecked

- Tokens
  * Token Type: JwtRS256
  * Auth methods:  check both "Client Creds" and "Resource Owner"

- Scope
  * "User must confirm authorization request" must be left unchecked
  * Add a new scope (e.g., testscope)
  * Allowed REST APIs must include the followings:
      - secrets/
      - privilegeddata/

- Permissions
  * grant "Run" permission to the test user

Test configuration file

You need to set up a JSON that describes the test configuration.  Here is an example
of the file:

 {
	"TenantURL": "tenant.my-centrify.net",
	"PASuser": { "Username" : "<specify the test PAS user here>",
		     "Password" : "<specify the password here>" },
	"AppID": "<name of web application setup above>",
	"Scope": "<scope setup above>"
 }

You can pass the configuration information in the go test command.  An example is <br>
`go test -args -config=/tmp/gotest.json`

*/
package secret
