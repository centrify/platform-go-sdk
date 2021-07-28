package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/centrify/cloud-golang-sdk/oauth"
	"github.com/centrify/platform-go-sdk/dmc"
	"github.com/centrify/platform-go-sdk/secret"
	"golang.org/x/term"
)

// structure to determine the operation selected
type actions struct {
	create       bool
	createFolder bool
	delete       bool
	get          bool
	getMetaData  bool
	list         bool
	modify       bool
}

type operation int

const (
	create operation = iota
	createFolder
	delete
	get
	getMetaData
	list
	modify
)

// Parameters defines the configuration parameters
type Parameters struct {
	// Configuration file
	ConfigFile string
	// URL where the secret is stored
	ServerPath string `json:"server"`
	// server type: pas, tss or dsv.  Only pas is supported in current version
	ServerType string `json:"servertype"`
	// whether to use DMC or not
	UseDMC bool `json:"useDMC"`
	// OAuth scope to use in REST calls
	Scope string `json:"scope"`

	// Other authentication related information when DMC is not used
	// Application ID
	AppID string `json:"appid"`
	// Client ID.
	ClientID string `json:"clientid"`
	// Client Secret.
	ClientSecret string `json:"clientsecret"`
	// username
	Username string `json:"user"`
	// password
	Password string `json:"password"`

	// Path to secret object/folder
	SecretPath string `json:"name"`
	// optional description of secret
	Description string `json:"description"`

	// the following parameters are usually specified in the command line

	// secret value for text secret
	TextValue string `json:"text"`
	// Keyvalue secret value stored in a JSON file
	JSONDataFile string `json:"jsonfile"`
	// Keyvalue secret value stored in a JSON string
	JSONString string `json:"jsonstring"`

	// These parameters are derived from other parameters and not specified in the
	// command line or in the configuration file.
	// type of secret operation
	Operation operation
	// Type of secret.  Must be one of "text" or "keyvalue"
	SecretType string
	KVSecret   map[string]string // content of keyvalue pair secret

	// These parameters are related to HTTP operations
	// whether to enable debug messages or not
	Debug bool `json:"debug"`
	// specify a different user agent in HTTP header
	UserAgent string `json:"useragent"`
	// specify extra HTTP headers in a comma separated list
	ExtraHeaders    string            `json:"headers"`
	ExtraHeadersMap map[string]string // extra headers stored as string map
	// whether to log REST API calls
	Log bool `json:"log"`
}

var errUsage error = errors.New("Usage error")

const usageConfig = `config file in JSON.  You can specify server, name, user, password, clientid, 
clientsecret, appID or scope. If a parameter is explicitly specified, it overrides the value in the file`
const usageText = `value of text secret to create/modify.  Must be specified when type is "text"`
const usageJSONFile = `JSON file that contains the keyvalue secret value to create/modify.  
Either JsonFile or JsonString must be specified when creating/modifying a keyvalue secret`
const usageJSONString = `A JSON string that specifies the keyvalue secret value to create/modify.
Either JsonFile or JsonString must be specified when creating/modifying a keyvalue secret`
const usageHeaders = `Specify extra HTTP headers as a comma-separated list.  Each header is specified as <name>:<value>.
Comma (,) and colon (:) are not allowed as part of the header name or value. 
Example: "X-TZOFF:480, X-Special:Marker`
const usageServerType = "Server type: pas for Centrify PAS"

// loadConfigFromFile loads the configuration parameters from a json file
func loadConfigFromFile(path string, result *Parameters) error {
	filebuf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error in reading file %s: %v\n", path, err)
		return err // error in read
	}
	err = json.Unmarshal(filebuf, result)
	if err != nil {
		fmt.Printf("Error in unmarshaling input file %s: %v\n", path, err)
		return err
	}
	return nil
}

// getBasicConfiguration gets the configuration based on command line parameters
func getConfiguration() (*Parameters, error) {

	// setup all the command line parameters
	cfgFilePtr := flag.String("config", "", usageConfig)

	// setup to read command line parameters into optCli so we can merge later
	cliOpt := &Parameters{}

	flag.StringVar(&cliOpt.ServerPath, "server", "", "Tenant URL where secret is stored")
	flag.StringVar(&cliOpt.ServerType, "servertype", "", usageServerType)
	flag.BoolVar(&cliOpt.UseDMC, "useDMC", false, "Use DMC. Note: It cannot be overridden if it is set to true in the config file.")
	flag.StringVar(&cliOpt.Scope, "scope", "", "scope")
	flag.StringVar(&cliOpt.AppID, "appid", "", "application ID")
	flag.StringVar(&cliOpt.Username, "user", "", "username")
	flag.StringVar(&cliOpt.Password, "password", "", "password")
	flag.StringVar(&cliOpt.ClientID, "clientid", "", "clientID")
	flag.StringVar(&cliOpt.ClientSecret, "clientsecret", "", "client Secret")
	flag.StringVar(&cliOpt.SecretPath, "name", "", "Path of secret")
	flag.StringVar(&cliOpt.Description, "description", "", "optional description of secret")
	flag.StringVar(&cliOpt.TextValue, "text", "", usageText)
	flag.StringVar(&cliOpt.JSONDataFile, "jsonfile", "", usageJSONFile)
	flag.StringVar(&cliOpt.JSONString, "jsonstring", "", usageJSONString)
	flag.BoolVar(&cliOpt.Debug, "debug", false, "Enable debug messages")
	flag.StringVar(&cliOpt.UserAgent, "useragent", "", "specify a different user agent in HTTP header")
	flag.StringVar(&cliOpt.ExtraHeaders, "headers", "", usageHeaders)
	flag.BoolVar(&cliOpt.Log, "log", false, "whether to log REST API call")

	// operation switch
	action := &actions{}

	flag.BoolVar(&action.create, "create", false, "create text or keyvalue secret")
	flag.BoolVar(&action.createFolder, "createfolder", false, "create secret folder")
	flag.BoolVar(&action.delete, "delete", false, "delete secret object/folder")
	flag.BoolVar(&action.get, "get", false, "get secret value")
	flag.BoolVar(&action.getMetaData, "getmetadata", false, "get secret metadata")
	flag.BoolVar(&action.list, "list", false, "list folder contents")
	flag.BoolVar(&action.modify, "modify", false, "modify secret")

	flag.Parse()

	options := new(Parameters)
	if *cfgFilePtr != "" {
		// config file specified, try to load it
		err := loadConfigFromFile(*cfgFilePtr, options)
		if err != nil {
			return nil, err
		}
		options.ConfigFile = *cfgFilePtr
	}

	mergeParameters(cliOpt, options)

	if !checkRequiredParameters(options) {
		return nil, errUsage
	}

	if !checkOperationSelection(options, action) {
		return nil, errUsage
	}

	if !checkCredSpecified(options) {
		return nil, errUsage
	}

	if !checkOptionalParameters(options) {
		return nil, errUsage
	}

	if !parseExtraHeaders(options) {
		return nil, errUsage
	}
	return options, nil
}

// getAccessToken returns the Oauth access token for the user
func getAccessToken(cfg *Parameters) (string, error) {
	if cfg.UseDMC {
		// get DMC Token
		token, err := dmc.GetDMCToken(cfg.Scope)
		if err != nil {
			return "", err
		}
		return token, nil
	}
	// get oauth token for user
	oauthClient, err := oauth.GetNewConfidentialClient("https://"+cfg.ServerPath, cfg.Username, cfg.Password, nil)
	if err != nil {
		fmt.Printf("Error in getting Oauth client: %v", err)
		return "", err
	}

	oauthToken, oauthError, err := oauthClient.ClientCredentials(cfg.AppID, cfg.Scope)
	if err != nil {
		fmt.Printf("Error in sending authentication request to server: %v", err)
		return "", err
	}

	if oauthError != nil {
		fmt.Printf("Authentication error: %v.  Description: %v\n", oauthError.Error, oauthError.Description)
		return "", errors.New(oauthError.Error)
	}
	return oauthToken.AccessToken, nil
}

// mergePerameters checks if any config parameter is specified in the command line, and use
// it to override the one specified in the config file.
// Note: If useDMC is set to true in the config file, it cannot be overridden.
func mergeParameters(cliOpt *Parameters, cfgOpt *Parameters) {
	if cliOpt.ServerPath != "" {
		cfgOpt.ServerPath = cliOpt.ServerPath
	}
	if cliOpt.ServerType != "" {
		cfgOpt.ServerType = cliOpt.ServerType
	}
	if cliOpt.UseDMC {
		cfgOpt.UseDMC = cliOpt.UseDMC
	}
	if cliOpt.AppID != "" {
		cfgOpt.AppID = cliOpt.AppID
	}
	if cliOpt.ClientID != "" {
		cfgOpt.ClientID = cliOpt.ClientID
	}
	if cliOpt.ClientSecret != "" {
		cfgOpt.ClientSecret = cliOpt.ClientSecret
	}
	if cliOpt.Username != "" {
		cfgOpt.Username = cliOpt.Username
	}
	if cliOpt.Password != "" {
		cfgOpt.Password = cliOpt.Password
	}
	if cliOpt.SecretPath != "" {
		cfgOpt.SecretPath = cliOpt.SecretPath
	}
	if cliOpt.Description != "" {
		cfgOpt.Description = cliOpt.Description
	}
	if cliOpt.TextValue != "" {
		cfgOpt.TextValue = cliOpt.TextValue
	}
	if cliOpt.JSONDataFile != "" {
		cfgOpt.JSONDataFile = cliOpt.JSONDataFile
	}
	if cliOpt.JSONString != "" {
		cfgOpt.JSONString = cliOpt.JSONString
	}
	if cliOpt.UserAgent != "" {
		cfgOpt.UserAgent = cliOpt.UserAgent
	}
	if cliOpt.Debug {
		cfgOpt.Debug = true
	}
	if cliOpt.ExtraHeaders != "" {
		cfgOpt.ExtraHeaders = cliOpt.ExtraHeaders
	}
	if cliOpt.Log {
		cfgOpt.Log = cliOpt.Log
	}
	if cliOpt.Scope != "" {
		cfgOpt.Scope = cliOpt.Scope
	}
}

// checkOperationSelection verifies that one and only one operation is selected
func checkOperationSelection(options *Parameters, selAction *actions) bool {

	var selOperation operation
	var optCount int

	if selAction.create {
		selOperation = create
		optCount++
	}
	if selAction.createFolder {
		selOperation = createFolder
		optCount++
	}
	if selAction.delete {
		selOperation = delete
		optCount++
	}
	if selAction.get {
		selOperation = get
		optCount++
	}
	if selAction.getMetaData {
		selOperation = getMetaData
		optCount++
	}
	if selAction.list {
		selOperation = list
		optCount++
	}
	if selAction.modify {
		selOperation = modify
		optCount++
	}

	if optCount > 1 {
		fmt.Println("Can only specify one of -create, -createFolder, -delete, -get, -getMetaData, -list or -modify")
		return false
	}
	if optCount == 0 {
		fmt.Println("Must specify one of -create, -createFolder, -delete, -get, -getMetaData, -list or -modify")
		return false
	}
	options.Operation = selOperation
	return true
}

// check credential requirement for PAS
func checkPasCred(options *Parameters) bool {
	if options.Scope == "" {
		fmt.Println("Must specify scope using -scope")
		return false
	}
	if !options.UseDMC {
		// DMC is not used, must specify username, password and AppID
		if options.Username == "" {
			fmt.Println("must specify user using -user")
			return false
		}
		if options.AppID == "" {
			fmt.Println("must specify application ID using -appid")
			return false
		}
		if options.Password == "" {
			// get password
			fmt.Print("Enter password: ")
			pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Printf("Error in reading password: %v\n", err)
				return false
			}
			options.Password = string(pwdBytes)
		}
	}
	return true
}

// checkCredSpecified checks if all information required to authenticate the user is specified
func checkCredSpecified(options *Parameters) bool {
	if options.ServerType == secret.ServerPAS {
		return checkPasCred(options)
	}

	// Note:  checkRequiredParameters already check whether server type is correct
	// Need to add additional credential check functions for server types such as DSV/TSS when support
	// is added
	fmt.Printf("%s is not supported\n", options.ServerType)
	return true
}

// checkRequiredParameters verify that all required parameters are present.
// this includes ServerPath and SecretPath
func checkRequiredParameters(options *Parameters) bool {
	if options.ServerPath == "" {
		fmt.Println("must specify tenant where secret is stored using -server")
		return false
	}
	if options.SecretPath == "" {
		fmt.Println("must specify secret path using -name")
		return false
	}
	if options.ServerType == "" {
		fmt.Println("must specify server type using -servertype")
		return false
	}
	options.ServerType = strings.TrimSpace(strings.ToLower(options.ServerType))
	if options.ServerType != secret.ServerPAS {
		fmt.Println("must specify \"pas\" as servertype")
	}
	return true
}

func checkOptionalParameters(options *Parameters) bool {

	if options.Operation != create && options.Operation != modify {
		// no need to check additional parameters
		return true
	}

	// for create and modify:
	// 1. Only one and only one of TextValue, JSONDataFile or JSONString must be specified

	var nSources int
	var secretType string
	var fromFile bool

	if options.TextValue != "" {
		nSources++
		secretType = secret.SecretTypeText
	}
	if options.JSONDataFile != "" {
		nSources++
		secretType = secret.SecretTypeKV
		fromFile = true
	}
	if options.JSONString != "" {
		nSources++
		secretType = secret.SecretTypeKV
	}
	if nSources != 1 {
		fmt.Printf("Must specify one and only one of -text, -jsonfile or -jsonstring")
		return false
	}
	options.SecretType = secretType

	if secretType == secret.SecretTypeKV {
		if fromFile {
			// verify that the file can be read and marshal into a map[string]string object
			filebuf, err := ioutil.ReadFile(options.JSONDataFile)
			if err != nil {
				fmt.Printf("Error in reading file %s: %v\n", options.JSONDataFile, err)
				return false
			}
			err = json.Unmarshal(filebuf, &options.KVSecret)
			if err != nil {
				fmt.Printf("Error in unmarshaling input data: %v\n", err)
				return false
			}
		} else {
			err := json.Unmarshal([]byte(options.JSONString), &options.KVSecret)
			if err != nil {
				fmt.Printf("Error in unmarshaling data: %v\n", err)
				return false
			}
		}
	}
	return true
}

// parseExtraHeaders parses the user specified comma separated list into a string map
func parseExtraHeaders(options *Parameters) bool {
	if options.ExtraHeaders == "" {
		options.ExtraHeadersMap = nil
		return true
	}
	hdrs := strings.Split(options.ExtraHeaders, ",")
	options.ExtraHeadersMap = make(map[string]string)
	for _, hdr := range hdrs {
		// split the string into name and value
		pair := strings.Split(hdr, ":")
		if len(pair) != 2 {
			fmt.Printf("[%s] is not a valid header as it is not a colon-separated string", hdr)
			return false
		}
		options.ExtraHeadersMap[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return true
}

func (op operation) String() string {

	var names = [...]string{
		"create",
		"createfolder",
		"delete",
		"get",
		"getmetadata",
		"list",
		"modify",
	}
	if op >= create && op <= modify {
		return names[op]
	}
	return "unknown"

}
