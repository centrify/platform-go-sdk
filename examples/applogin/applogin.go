package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/centrify/platform-go-sdk/vault"
)

type parameters struct {
	scope    string // DMC scope
	vaultURL string
}

var errUsage error = errors.New("Usage error")

// getParameters gets parsed command line parameters
func getParameters() (*parameters, error) {

	opt := &parameters{}
	flag.StringVar(&opt.scope, "scope", "", "DMC scope")
	flag.StringVar(&opt.vaultURL, "url", "", "HashiCorp Vault url")
	flag.Parse()

	if opt.scope == "" {
		fmt.Println("scope must be specified")
		return nil, errUsage
	}
	if opt.vaultURL == "" {
		fmt.Println("url must be specified")
		return nil, errUsage
	}
	return opt, nil
}

func main() {
	parameters, err := getParameters()
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
	token, err := vault.GetHashiVaultToken(parameters.scope, parameters.vaultURL)

	if err != nil {
		fmt.Printf("Error in getting HashiCorp Vault token: %v\n", err)
		os.Exit(2)
	}
	// success, just print token in stdout
	fmt.Print(token)
	os.Exit(0) // so caller can check exit status
}
