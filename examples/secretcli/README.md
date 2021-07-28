# secretcli - A sample program on accessing secrets using the secret package

## Usage

The following command line parameters are supported:
```
  -appid string
    	application ID
  -config string
    	config file in JSON.  You can specify server, name, user, password, appID or scope.
    	If a parameter is explicitly specified, it overrides the value in the file
  -create
    	create text or keyvalue secret
  -createfolder
    	create secret folder
  -debug
    	Enable debug messages
  -delete
    	delete secret object/folder
  -description string
    	optional description of secret
  -get
    	get secret value
  -getmetadata
    	get secret metadata
  -headers string
    	Specify extra HTTP headers as a comma-separated list.  Each header is specified as <name>:<value>.
    	Comma (,) and colon (:) are not allowed as part of the header name or value. 
    	Example: "X-TZOFF:480, X-Special:Marker
  -jsonfile string
    	JSON file that contains the keyvalue secret value to create/modify.  
    	Either JsonFile or JsonString must be specified when creating/modifying a keyvalue secret
  -jsonstring string
    	A JSON string that specifies the keyvalue secret value to create/modify.
    	Either JsonFile or JsonString must be specified when creating/modifying a keyvalue secret
  -list
    	list folder contents
  -log
    	whether to log REST API call
  -modify
    	modify secret
  -name string
    	Path of secret
  -password string
    	password
  -scope string
    	scope
  -server string
    	Tenant UTL where secret is stored
  -servertype string
    	Server type: pas for Centrify PAS
  -text string
    	value of text secret to create/modify.  Must be specified when type is "text"
  -useDMC
    	Use DMC. Note: It cannot be overridden if it is set to true in the config file.
  -user string
    	username
  -useragent string
    	specify a different user agent in HTTP header
```
## Use a JSON file to store commonly used parameters

You can store these commonly used parameters in a JSON file and specify the JSON file using the -config parameter.
- appid
- debug
- description
- headers
- jsonfile
- jsonstring
- log
- name
- password (not recommended)
- scope
- server
- servertype
- text
- useDMC
- user
- useragent

Here is an example of a jsonfile
```
{
	"appid": "testsdk",
	"server": "my-tenant.centrify.com",
	"user": "developer@acme.com",
	"scope": "all",
	"servertype": "pas"
}
```

## Use Delegated Machine Credential to access secret

To avoid saving password/credential information in your scripts/programs, you can use Delegated Machine Credentials to access secrets.

### Pre-requisites
1. You need to install Centrify Client on the machine, 
2. The machine must login to your tenant with DMC feature enabled.  For example:
```
cenroll -t ${TENANT_URL} -c ${CODE} -F dmc -d secret=secrets\$ -d secret=secrets/.\* -d secret=privilegeddata/.\*
```
Note that ${TENANT_URL} is the environment variable for the tenant URL; and ${CODE} is the environment variable for the enrollment code.

Alternatively, you can set up the same DMC scope by editing the Client Profile of the Centrify Client in the admin portal.

3. You need to run this program as root.  Only root users can acquire Delegated Machine Credential.

4. Since the secret is accessed using the machine account, make sure that the machine account is granted the correct permissions to the secrets/folders.

### Sample JSON configuration file when using Delegated Machine Credential
```
{
	"useDMC": true,
	"server": "my-tenant.centrify.com",
	"scope": "secret",
	"servertype": "pas"
}
```
## Exit status

| Status | Errors |
| ------ | ------ |
| EPERM(1) | ErrSecretTypeNotSupported: Specified secret type is not supported. |
|  | ErrCannotModifySecretType: Cannot modify secret type.|
| | ErrCannotModifySecretFolder: Cannot modify secret folder. |
| ENOENT(2) | ErrFolderNotFound: Secret folder does not exist. |
|| ErrSecretNotFound: Secret does not exist. |
| EACCES(13) | ErrNoCreatePermission: No permission to create secret/folder. |
| | ErrNoDeletePermission: No permission to delete secret/folder. |
| | ErrNoGetMetaDataPermission: No permission to get metadata information about secret/folder. |
| | ErrNoModifyPermission: No permission to modify secret. |
| | ErrNoRetrievePermission: No permission to retrieve secret. |
| EEXIST(17) | ErrExists: Secret alreay exists. |
|| ErrDeletedSecretExists: A mark-for-delete secret already exists in the same path. |
| ENOTDIR(20) | ErrNotSecretFolder: Path is not a secret folder. |
| EISDIR(21) | ErrNotSecretObject: Path is a secret folder. |
| EINVAL(22) | ErrBadPathName: Illegal secret path name. |
| | ErrBadServerType: Invalid server type. |
| ENOSYS(38) | ErrNotImplementedYet: Function not implemented yet. |
| ENOTEMPTY(39) | ErrFolderNotEmpty: Secret folder is not empty. |
| EPROTO(72) | ErrUnexecptedResponse: Unexpected response received. |
| 255 | Usage error. Error in command line parameters. |

## Examples
### Listing secrets in a folder
```
$ sudo ./secretcli -config ~/dmc.json -name folder1 -list
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1     list  map[] false   map[] false}
Listing contents of [folder1]
Number of items in folder: 4
ID: 1fd46425-49dd-4cb3-bbea-783dfb32ab68	Type: Folder	Name: folder3
ID: 90d07161-07df-4464-bc37-71899d7dc2be	Type: Folder	Name: textsecret
ID: 0cb524cc-2b97-4084-87ec-fd111fc588ac	Type: Text		Name: newsecrettext
ID: 0d59ddd6-7faf-4efc-9b87-be3033506597	Type: KeyValue	Name: bag-secret
```
### Getting values of a secret

A text secret:
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/newsecrettext -get
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/newsecrettext     get  map[] false   map[] false}
Getting secret from path [folder1/newsecrettext]
Secret is a text string. Value: [now i change it]
```
A keyvalue secret:
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/bag-secret -get
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/bag-secret     get  map[] false   map[] false}
Getting secret from path [folder1/bag-secret]
Secret is key value pair collection:
Key: key3	Value:third_value
Key: key with space	Value:value with space
```
### Create a text secret
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-is-fun -create -text "That's all folks"'!'
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-is-fun  That's all folks!   create text map[] false   map[] false}
Creating secret of type text in path [folder1/secret-is-fun]
Secret created. ID: 428f658f-f5aa-4be7-83d0-e9e97ddd5c0e

$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-is-fun -get
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-is-fun     get  map[] false   map[] false}
Getting secret from path [folder1/secret-is-fun]
Secret is a text string. Value: [That's all folks!]
```
### Create a keyvalue secret
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-keyvalue -create -jsonstring "{\"foo\":\"bar\", \"hello\":\"world\"}"
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-keyvalue    {"foo":"bar", "hello":"world"} create keyvalue map[foo:bar hello:world] false   map[] false}
Creating secret of type keyvalue in path [folder1/secret-keyvalue]
Secret created. ID: 0b9760fe-2ae4-4029-b1ee-2fb3a255873c

$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-keyvalue -get
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-keyvalue     get  map[] false   map[] false}
Getting secret from path [folder1/secret-keyvalue]
Secret is key value pair collection:
Key: hello	Value:world
Key: foo	Value:bar
```
### Modify a secret
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-keyvalue -modify -jsonstring "{\"bar\":\"foo\", \"world\":\"hello\"}"
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-keyvalue    {"bar":"foo", "world":"hello"} modify keyvalue map[bar:foo world:hello] false   map[] false}
Modifying secret of type keyvalue in path [folder1/secret-keyvalue]
Secret modified. ID: 0b9760fe-2ae4-4029-b1ee-2fb3a255873c

$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-keyvalue -get
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-keyvalue     get  map[] false   map[] false}
Getting secret from path [folder1/secret-keyvalue]
Secret is key value pair collection:
Key: world	Value:hello
Key: bar	Value:foo
```
### Delete a secret
```
$ sudo ./secretcli -config ~/dmc.json -name folder1/secret-is-fun -delete
Configuration: {/home/user/dmc.json my-tenant.centrify.com pas true secret    folder1/secret-is-fun     delete  map[] false   map[] false}
Deleting secret in path [folder1/secret-is-fun]
Secret deleted
```
