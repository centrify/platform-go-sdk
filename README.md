[![Build Status](https://github.com/centrify/platform-go-sdk/workflows/build/badge.svg)](https://github.com/centrify/platform-go-sdk/actions)

# platform-go-sdk
Public Go packages for using Centrify Platform.

Subdirectories:

- DMC: Allow applications to acquire Delegated Machine Credentials token
- Examples: Sample programs that demonstrate how some packages can be used
  * AppLogin: A sample of how to get a HashiCorp Vault token.
  * DMC: An example on how to get DMC tokens.
  * SecretCLI:  A CLI program that can be used to access secrets.
- OAuthhelpers: Export a public method that retrieves an OAuth token using Resource Owner grant request. This can only be used by Centrify Vault software. Contact ThycoticCentrify support if you need to use this API.
- Secret: Allow applications to create/read/update/delete PAS secrets.
- TestUtils: Support functions that can be used by go tests.
- Utils: Miscellaneous methods for getting information about current system.
- Vault: HashiCorp Vault related functions.


## License

See [LICENSE](LICENSE)
