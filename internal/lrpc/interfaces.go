/*
Package lrpc defines various interfaces and types used in communication with Centrify Client

Porting to different operating systems

The current implementation is implemented as tested in Linux.  System dependent code are implemented in the files peer_linux.go and network_linux.go; and you need to provide equivalent functionalities for other platforms.

File session_linux.go

 This file implements all OS dependent functions defined in the SessionCtxt interface.  You need to implement the factory function
 NewSessionCtxt() and the following functions in SessionCtxt interface:

	// IsPrivileged() returns true if the peer process runs as privileged process.
	// In *nix, this means the effective UID is root (0)
	// In windows, this means the owner of process has local admin right
	IsPrivileged() (bool, error)

	// GetProcessID() returns the process ID of the peer process
	GetProcessID() (int32, error)

	// GetProgram() returns the name of program executable that the peer process is running with
	GetProgram() (string, error)

 Note that SessionCtxtBase is a helper struct that provides OS independent support for SessionCtxt.  More importantly, it implements the
 Get() and Set() methods for managing session related information.

File network_linux.go

 // This file implements the OS-dependent code for listening and connecting to the service endpoint.  You need to implement the following functions:
 //
 // StartListener is OS-dependent function to establish an internal listener.
 // Input parameter:
 // endpoint (string):          The name of the endpoint
 // acl (AccessControlLevel):   The access control level to start listening
 //
 // Return values:
 // net.Listener:	network listener object established.
 // error
 //
 // Note:
 // This function should remove the endpoint if it already exists.

 func StartListener(endpoint string, acl AccessControlLevel) (net.Listener, error)

 // ConnectToServer is OS-dependent code for connecting to an internal endpoint.
 // It is used by a RPC client.
 //
 // Input parameters:
 // endpoint (string):          The endpoint to connect to.
 // acl (AccessControlLevel):   The ACL for establishing the listener
 //
 // Return values:
 // net.Conn:   The established network connection.
 // error:		err

 func ConnectToServer(endpoint string) (net.Conn, error)

Writing a RPC session server

You need to consider the followings when you implement a LRPC session server that serves all the RPC messages to an endpoint.

1. The process model.   This depends on the nature of the messages/services.  For most cases, it can be done simply by having a goroutine to serve all the messages for a LRPC connection.   However, for NSS-related messages in Linux, it may be better to use a pre-configured pool of goroutines to serve all the requests as this will throttle such requests.   This may or may not be necessary.  Need to analyze performance test results first.

2. Message handling.   For Windows, we may not need to protect the messages from snooping, and there may not be any need to implement sequence numbering and signing of request/request.  However, we need to do this in *nix for security reasons.

This package currently provide a function that creates a RPC session server:
  NewLrpc2SessionServer - creates a session server that uses a goroutine to serve all the messages for each LRPC connection.


Messages

A message can be identified by message ID.  You can choose any message number that you like for your messages.  The only limitations are:

1. The mesasge ID MUST be unique for each endpoint.

2. Message ID must convert correctly into an uint16 value.

You need to register all the messages for an endpoint before you invoke the Start() method.

Message handler

You need to implement a handler for each message as type CommandHandler.  See description on type CommandHandler for details.

Writing application that needs LRPC service

1. Create a client connection by calling NewLrpc2ClientSession().

2. Connect to the client using Connect()

3. Set up the arguments for the LRPC call as an array of interface{}.

4. Call lrpc.DoRequest() (if a response is expected) or lrpc.DoAsyncRequest() (if no response is expected).

5. When lrpc.DoRequest() is called, the results are passed back as an array of interface{}.  Process the results.


Steps in writing a session server for a specific endpoint

1. Call NewLrpc2SessionServer() (or provide your own function to create a session handler that implements the SessionServer interface).

2. Register the messages for the handler by calling one or more of the followings on the created message handler:
	RegisterMsgsByID()
	RegisterMsgByID()

3. Call the Start() method of the handler to start the service.

4. Call the Stop() method of the handler to stop the service.

5. Call the Wait() method of the handler to wait for the service to be completely shutdown.

6. Implement a handler function for each message.

7. If you need to save information for use for other messages sharing the same lrpc session, you can use the functions Get() and Set()
   in the associated SessionCtxt interface.
*/
package lrpc

import (
	"unsafe"
)

/*
CommandHandler is the function prototype of a handler that handles a LRPC message.

 Input parameters:
    ctxt: The context of the session.  It implements the SessionCtxt interface.
          The handler can call methods in this interface to get information about the requester,
          as wells managing information across multiple messages sharing the same LRPC session.
    args: An array of interfaces that contains the arguments for the command.

 Output:
	result:  An array of interfaces that contains the results for the command.

Notes for CommandHandler:

1. Each element of the arguments and results can have its own type.  The command handler can validate the types of each argument to verify it is what it expects. Similarly, the handler is responsible for using the correct type of each element in the result array.

2. The command handler MUST return nil if the client does not expect any response.

3. String type argument/result MUST NOT contain the null character ('\x00')

*/
type CommandHandler func(ctxt SessionCtxt, args []interface{}) []interface{}

/*
MessageClient defines all the required functions that a message client must implement.

Notes:

1. The implementation needs to ensure that ReadResponse() will return the correct response for the previous WriteRequest()/WriteNamedRequest().

2. For simplicity, the implementations assume that only one thread/goroutine can run WriteRequest()/ReadRequest().   If the object is shared, the caller MUST synchronize access to the object.

3. Legacy LRPC2 message protocol DOES NOT support calling messages by name.  In this case, the first parameter in WriteRequest parameter must be uint32
*/
type MessageClient interface {
	//	Connect initiates a connection to the remote endpoint (specified when the MessageClient is created)
	Connect() error

	// WriteRequest sends a request message to the server.
	// Parameters:
	// cmd:	Command to send (only support uint32 and string type)
	// args:	Command arguments
	WriteRequest(cmd interface{}, args []interface{}) error

	// ReadResponse:	read the response
	// Return values:
	// 	results:	results in the response
	// 	err:		error
	ReadResponse() (results []interface{}, err error)

	// IsNamedMessagesSupported() returns a bool to indicate whether calling messages by name is supported.
	// Return value:
	// 	bool:	whether calling messages by name is supported
	IsNamedMessagesSupported() bool

	// Close() closes the client connection to the server
	// Return value:
	// 	err:	error
	Close() error
}

/*
SessionCtxt specifies the interface for getting information about the the current LRPC session.
Note that the implementation is very likely to be system dependent.

Each implemenation of MessageServer type must implement the function GetSessionCtxt() that returns an object
that supports this interface
*/
type SessionCtxt interface {

	// IsPrivileged() returns true if the peer process runs as privileged process.
	// In *nix, this means the effective UID is root (0)
	IsPrivileged() (bool, error)

	// GetProcessID() returns the process ID of the peer process
	GetProcessID() (int32, error)

	// GetProgram() returns the name of program executable that the peer process is running with
	GetProgram() (string, error)

	// GetCallerUserID() returns user ID of the caller.
	// For Linux, it is the string representation of the UID
	// For windows, it is the SID in string format
	GetCallerUserID() (string, error)

	// Get() gets an attribute associated with the object.  Both key and returned value are
	// application dependent and are defined by the command handlers
	Get(key string) interface{}

	// Set() sets an attribute associated with the object.  The key and val are application dependent.
	Set(key string, val interface{})
}

/*
MessageServer defines all the required functions that a message server must implement.

Notes:

1. 'msg' is used to associate a response to the original request.   It should be treated as opaque for the callers to ReadRequest/WriteResponse.

2. For simplicity, the implemenations do not serialize access to ReadRequest()/WriteResponse().  It is up to the caller to synchronize access if necessary.
*/
type MessageServer interface {
	// Accept() accepts a client connection and returns a new object that represents the accepted connection and supports the MessageServer interface.
	// 	Return values:
	// 	newServer:	New object that supports the MessageServer interface
	// 	err:		error
	Accept() (MessageServer, error)

	// GetSessionCtxt() returns the context for the LRPC session
	// 	Return values:
	// 	peer:	object that supports the SessionCtxt interface
	// 	err:	error
	GetSessionCtxt() (SessionCtxt, error)

	// ReadRequest() reads a request.
	// 	Return values:
	// 	msg:		opaque context for the message.  Need to be passed to WriteResponse()
	// 	command:	command ID
	// 	args:		arguments for the command
	// 	err:		error
	ReadRequest() (msg interface{}, command interface{}, args []interface{}, err error)

	// WriteResponse() writes the resposne to a command.
	// 	Parameters:
	// 	msg:		message context as received in ReadRequest
	//  command:		command ID
	// 	results:	results for the command.
	// 	Return value:
	// 	err:		error
	WriteResponse(msg interface{}, command interface{}, results []interface{}) error

	// IsNamedMessagesSupported() returns a bool to indicate whether calling messages by name is supported.
	// 	Return values:
	// 	bool:	whether calling messages by name is supported
	IsNamedMessagesSupported() bool

	// Close() closes the network connection associated with the current server connection
	Close() error
}

/*
AccessControlLevel provides Access Control Level for creating listener. The is platform specific that how to use this data
*/
type AccessControlLevel interface{}

// Constants to represent the various data types in message
//
// Notes:
//
// 1. The values for constants msgEnd through msgDataTypeProtectedBlob MUST NOT
// be changed.  They need to match the corresponding values in used by Centrify Client.
//
// 2. There should be no need for any applications to know about these constants.
const (
	msgEnd                   byte = iota // Message end..no data
	msgDataTypeBool                      // type bool
	msgDataTypeInt32                     // type int32
	msgDataTypeUint32                    // type uint32
	msgDataTypeString                    // type string
	msgDataTypePassword                  // deprecated...
	msgDataTypeBlob                      // type []byte
	msgDataTypeStringSet                 // type []string
	msgDataTypeKeyValueSet               // type map[string]string
	msgDataTypeProtectedBlob             // Not supported
	// the following types ARE NOT supported in C/C++ code yet for LRPC...
	// Do not use them unless you are sure that the message does not communicate with C/C++ code
	msgDataTypeByte   // type byte/uint8
	msgDataTypeUint64 // type uint64
	msgDataTypeInt64  // type int64
	msgDataTypeInt    // type int
	// The following message types are special
	msgDataTypeNil // nil object.  Implemented as special string of length -1 in LRPC2,
	// but should have a different tag for other implementations
	// Special note:
	// LRPC2 encodes []uint32 as a sequence of uint32 data entities.  However, they will be decoded as individual uint32 data objects.
	// Be careful when passing it as arguments or results
)

/*
SessionServer specifies the methods that a session server must implement.
*/
type SessionServer interface {
	// note: the following 4 functions are pre-implemented in the base struct baseSvc
	RegisterMsgsByID(map[uint16]interface{}) error
	RegisterMsgByID(uint16, interface{}) error

	Start() error
	Stop() error
	Wait() error
}

var sizeOfInt uintptr

// determine size of int
func init() {
	var i int
	sizeOfInt = unsafe.Sizeof(i)
}
