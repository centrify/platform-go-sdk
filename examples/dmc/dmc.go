package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/centrify/platform-go-sdk/dmc"
)

type parameters struct {
	scope string // DMC scope
}

var errUsage error = errors.New("Usage error")

// getParameters gets parsed command line parameters
func getParameters() (*parameters, error) {

	opt := &parameters{}
	flag.StringVar(&opt.scope, "scope", "", "DMC scope")
	flag.Parse()

	if opt.scope == "" {
		fmt.Println("scope must be specified")
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

	token, err := dmc.GetDMCToken(parameters.scope)
	if err != nil {
		fmt.Printf("Error in getting DMC token: %v\n", err)
		os.Exit(2)
	}
	// success, just print token in stdout
	fmt.Print(token)
	os.Exit(0) // so caller can check exit status
}
