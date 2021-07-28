package lrpc

import (
	"encoding/binary"
	"errors"
	"time"
)

const (
	//
	//
	// lrpc2Nack represents a NACK message
	lrpc2Nack = uint32(0)
	// lrpc2Ack represents an ACK message
	lrpc2Ack = uint32(1)

	// We include a magic number in all messages. Four byte value.
	lrpc2MagicNum = uint32(0xABCD8012)

	//
	// The maximum LRPC2 message data length. Currently only applicable to
	// LRPC2 requests. There is no limit on LRPC2 responses.
	//
	// The maximum thread number would be less than 20. Increase this value
	// when you think necessary.
	//
	lrpc2MaxMsgLen = uint32(1024 * 1024)

	// Handshake request includes LRPC2 protocol version
	lrpc2HandshakeRequestSize = 4 // bytes

	// Supported LRPC2 versions
	lrpc2Version4 = uint32(4)

	// LRPC2 version 4 settings
	lrpc2HandshakeReplyV4Size = 8          // bytes
	lrpc2HeaderLengthV4       = uint16(34) // bytes
)

// Message ID
const (
	LrpcMsgIDClientInfo            = 119
	Lrpc2MsgIDAdminClientGetToken  = 1500
	Lrpc2MsgIDGetPublicKey         = 1501
	Lrpc2MsgIDGetResourceOwnerCred = 1502
	Lrpc2MsgGetHashicorpVaultToken = 1503
)

// lrpc2ByteOrder defines the byte order used in LRPC messages
var lrpc2ByteOrder = binary.LittleEndian

// LRPC2 errors
var (
	ErrLrpc2BadMagicNum       = errors.New("Bad magic number")
	ErrLrpc2BadHeaderLen      = errors.New("Bad header length")
	ErrLrpc2MsgVerMismatch    = errors.New("Message version mismatched")
	ErrLrpc2BadMsgLen         = errors.New("Bad message length")
	ErrLrpc2ProcessIDMismatch = errors.New("Process ID mismatched")
	ErrLrpc2SeqNumMismatch    = errors.New("Sequence number mismatched")
	ErrLrpc2NameNotSupported  = errors.New("Send command by name is not supported")
	ErrLrpc2MsgTooLong        = errors.New("Message length exceeds LRPC2 limit")
	ErrLrpc2PELocalUser       = errors.New("Local User is not supported")
	ErrLrpc2TypeNotSupported  = errors.New("Data type not supported")
)

//
// Message status while Processing LRPC2 Message, these have same status
// name as 'C' part.
//
var (
	ErrMsgUnexpectedData   error = errors.New("Lrpc message contain unexpected data")
	ErrMsgIncorrectType    error = errors.New("Incorrect data type found in lrpc message")
	ErrMsgAtMsgEnd         error = errors.New("Already at end of lrpc message")
	ErrMsgNotComplete      error = errors.New("Lrpc message is incomplete")
	ErrMsgEncTypeNoSupp    error = errors.New("Enc type is not suported for encrypt/decrypt lrpc message")
	ErrMsgEncError         error = errors.New("Error while encrypt/decrypt protected lrpc message")
	ErrMsgInvalid          error = errors.New("Lrpc Message is invalid")
	ErrMsgBadHandshakeSize error = errors.New("LRPC2 handshake request size mismatched")
	ErrMsgBadVersion       error = errors.New("Unknown LRPC2 version")
	ErrMsgBadContext       error = errors.New("Incorrect response context")
)

// Common server errors
var (
	ErrLrpcServerCommandOutOfRange = errors.New("Command out of supported range")
	ErrLrpcServerNotConnected      = errors.New("Not connected")
	ErrLrpcServerAlreadyConnected  = errors.New("Already connected")
)

// timeout
const (
	lrpc2ConnectTimeout = 5 * time.Second
	lrpc2ReceiveTimeout = 300 * time.Second
	lrpc2SendTimeout    = 60 * time.Second
)
