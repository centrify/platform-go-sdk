package lrpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

// type lrpc2Client represents a LRPC2 client session.  It must implement the MessageClient interface
type lrpc2Client struct {
	lrpc2HeaderV4
	endpoint      string             // end point information
	conn          net.Conn           // network connection
	config        *lrpc2ClientConfig // lrpc2 client config
	maxMsgDataLen uint32             // maximum data size
	sessionPid    uint64             // session Pid
}
type lrpc2ClientConfig struct {
	// LRPC2 client-side timeout when establishing connection (handshake)
	ConnectTimeout time.Duration

	// LRPC2 client-side timeout when reading data from server reply
	ReceiveTimeout time.Duration

	// LRPC2 client-side timeout when sending request to server
	SendTimeout time.Duration
}

func newLrpc2ClientConfig() *lrpc2ClientConfig {
	return &lrpc2ClientConfig{
		ConnectTimeout: lrpc2ConnectTimeout,
		ReceiveTimeout: lrpc2ReceiveTimeout,
		SendTimeout:    lrpc2SendTimeout,
	}
}

// initLrpc2ClientSession init a new LRPC2 client session object.
func initLrpc2ClientSession(endpoint string, config *lrpc2ClientConfig) *lrpc2Client {
	c := new(lrpc2Client)
	c.magicNum = lrpc2MagicNum
	c.headerLen = lrpc2HeaderLengthV4
	c.version = lrpc2Version4
	c.endpoint = endpoint
	c.config = config

	return c
}

// NewLrpc2ClientSession create a new LRPC2 client session object.
func NewLrpc2ClientSession(endpoint string) MessageClient {
	config := newLrpc2ClientConfig()
	c := initLrpc2ClientSession(endpoint, config)
	if config == nil || c == nil {
		logging.Errorf("LRPC Client: Cannot init client session for endpoint %s", endpoint)
		return nil
	}

	// seed random number
	rand.Seed(time.Now().UnixNano())
	return c
}

func (c *lrpc2Client) IsNamedMessagesSupported() bool {
	return false
}

func (c *lrpc2Client) doHandShake() error {
	var err error

	// Note that the connection timeout is for the whole handshake. So no
	// need to set again for every read/write operation.
	now := time.Now()
	err = c.conn.SetDeadline(now.Add(c.config.ConnectTimeout))
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", c.conn, err)
		return err
	}

	logging.Tracef("LRPC2 client: Handshaking for LRPC2 connection [%p] with timeout %v...", c.conn, c.config.ConnectTimeout)

	var reqbytes = make([]byte, 0, lrpc2HandshakeRequestSize)
	var req = bytes.NewBuffer(reqbytes)

	err = binary.Write(req, lrpc2ByteOrder, uint32(lrpc2Version4))
	if err != nil {
		logging.Errorf("LRPC2 client: Internal error. Failed to write LRPC version information to buffer: %v", err)
		return err
	}

	_, err = c.conn.Write(req.Bytes())
	if err != nil {
		logging.Errorf("LRPC2 client: Internal error. Failed to send LRPC version information to server: %v", err)
		return err
	}

	logging.Trace("LRPC2 client: Sent handshake request, waiting for reply")

	bytesReply := make([]byte, lrpc2HandshakeReplyV4Size)
	_, err = c.conn.Read(bytesReply)
	if err != nil {
		logging.Debugf("LRPC2 client: Failed to receive handshake reply: %v", err)
		return err
	}

	var answer uint32
	var maxMsgDataLen uint32

	reply := bytes.NewBuffer(bytesReply)

	err = binary.Read(reply, lrpc2ByteOrder, &answer)
	if err != nil {
		logging.Debugf("LRPC2 client: Failed to read answer from handshake reply: %v", err)
		return err
	}

	if answer == lrpc2Nack {
		message := "LRPC2 handshake rejected"
		logging.Debugf("LRPC2 client: %s", message)
		return fmt.Errorf(message)
	}

	err = binary.Read(reply, lrpc2ByteOrder, &maxMsgDataLen)
	if err != nil {
		logging.Debugf("LRPC2 client: Failed to read mesasge body size limit from handshake reply: %v", err)
		return err
	}

	c.maxMsgDataLen = maxMsgDataLen
	c.sequenceNum = rand.Uint32()
	c.pid = uint64(os.Getpid())
	c.sessionPid = c.pid
	logging.Tracef("LRPC2 client handshake completed. MaxMsgDataLen: %d, seqNum: %d PID: %d",
		c.maxMsgDataLen, c.sequenceNum, c.pid)
	return nil
}

func (c *lrpc2Client) Connect() error {
	logging.Tracef("LRPC2 client: Connecting to LRPC server %s...", c.endpoint)

	conn, err := ConnectToServer(c.endpoint)
	if err != nil {
		logging.Debugf("LRPC2 client: cannot connect to %s: %v", c.endpoint, err)
		return err
	}

	c.conn = conn

	err = c.doHandShake()
	if err != nil {
		errClose := conn.Close()
		if errClose != nil {
			logging.Errorf("LRPC2 client: Failed to close LRPC connection [%p] after handshake failure: %v",
				conn, errClose)
		}
		return err
	}
	logging.Trace("LRPC2 client: Handshake completed successfully")
	return nil
}

func (c *lrpc2Client) WriteRequest(cmd interface{}, args []interface{}) error {
	var err error

	now := time.Now()
	err = c.conn.SetDeadline(now.Add(c.config.SendTimeout))
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", c.conn, err)
		return err
	}

	logging.Tracef("LRPC2 client: Sending request to LRPC2 connection [%p] with timeout %v...", c.conn, c.config.SendTimeout)

	defer c.setupNextRequest(c.sequenceNum + 1)

	// note that the cmd may be specified as an int, uint, uint16, uint32
	// make sure that it can fit into a uint16 value
	var iReq uint16

	switch cmd.(type) {
	case int:
		v := cmd.(int)
		iReq = uint16(v)
		if uint64(iReq) != uint64(v) {
			logging.Errorf("LRPC2 client: Command value %v of type %T does not fit into uint16.", cmd, cmd)
			return ErrLrpcServerCommandOutOfRange
		}

	case uint:
		v := cmd.(uint)
		iReq = uint16(v)
		if uint64(iReq) != uint64(v) {
			logging.Errorf("LRPC2 client: Command value %v of type %T does not fit into uint16.", cmd, cmd)
			return ErrLrpcServerCommandOutOfRange
		}

	case uint16:
		iReq = cmd.(uint16)

	case uint32:
		v := cmd.(uint32)
		iReq = uint16(v)
		if uint64(iReq) != uint64(v) {
			logging.Errorf("LRPC2 client: Command value %v of type %T does not fit into uint16.", cmd, cmd)
			return ErrLrpcServerCommandOutOfRange
		}

	default:
		// LRPC2 only support messages by ID
		logging.Errorf("LRPC2 client: Sending command by name [%v] (type %T) is not supported.", cmd, cmd)
		return ErrLrpc2NameNotSupported
	}

	// LRPC2 expects the command to be type uint16
	msgData, err := encode(iReq, args)
	if err != nil {
		return err
	}

	// check if size is supported
	if uint64(msgData.Len()) > uint64(c.maxMsgDataLen) {
		return ErrLrpc2MsgTooLong
	}

	// allocate a big byte array to store the full message
	c.msgDataLen = uint32(msgData.Len())
	bytesMsg := make([]byte, 0, c.HeaderLen()+int(c.msgDataLen))
	msg := bytes.NewBuffer(bytesMsg)

	// set up message header
	c.timestamp = uint64(time.Now().Unix())
	err = c.encodeHeader(msg)
	if err != nil {
		return err
	}

	// Add real data
	err = binary.Write(msg, lrpc2ByteOrder, msgData.Bytes())
	if err != nil {
		logging.Errorf("LRPC client: Cannot write message content to out buffer: %v", err)
		return err
	}

	// send data to remote
	_, err = c.conn.Write(msg.Bytes())
	logging.Tracef("LRPC client: message sent. sequence number: %d", c.sequenceNum)
	return err
}

// ReadResponse() reads the response for the request just sent....
func (c *lrpc2Client) ReadResponse() ([]interface{}, error) {
	var err error

	now := time.Now()
	err = c.conn.SetDeadline(now.Add(c.config.ReceiveTimeout))
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", c.conn, err)
		return nil, err
	}

	logging.Tracef("LRPC2 client: Receiving reply from LRPC2 connection [%p] with timeout %v...", c.conn, c.config.ReceiveTimeout)

	logging.Tracef("LRPC client: Entering ReadResponse for request ID: %d", c.sequenceNum)

	expectSeq := c.sequenceNum - 1

	bytesMsgHeader := make([]byte, c.HeaderLen())

	_, err = c.conn.Read(bytesMsgHeader)
	if err != nil {
		logging.Errorf("LRPC client: Error in reading response header: %v", err)
		return nil, err
	}

	msgHeader := bytes.NewBuffer(bytesMsgHeader)
	err = c.decodeHeader(msgHeader)
	if err != nil {
		logging.Errorf("LRPC client: Error in decoding response header: %v", err)
		return nil, err
	}

	// verify header
	err = c.verifyHeader()
	if err != nil {
		logging.Errorf("LRPC client: Error in verifying header: %v", err)
		return nil, err
	}

	// verify sequence number
	if c.sequenceNum != expectSeq {
		logging.Errorf("LRPC client: Expect response sequence number: %d, got %d", expectSeq, c.sequenceNum)
		return nil, ErrLrpc2SeqNumMismatch
	}

	logging.Tracef("LRPC client: Reading message data from LRPC connection [%p].  Expected size: %d bytes", c.conn, c.msgDataLen)

	bytesMsgData := make([]byte, c.msgDataLen)
	_, err = c.conn.Read(bytesMsgData)
	if err != nil {
		logging.Errorf("LRPC client: Error in reading message data from LRPC connection [%p]: %v", c.conn, err)
		return nil, err
	}

	msgData := bytes.NewBuffer(bytesMsgData)
	logging.Tracef("LRPC client: Received %d bytes of message data from LRPC connection [%p] ", c.msgDataLen, c.conn)

	_, rest, err := decode(msgData)
	if err != nil {
		logging.Errorf("LRPC client: Error in decoding response: %v", err)
		return nil, err
	}

	logging.Tracef("LRPC client: Return  %d values", len(rest))
	return rest, nil

}

// setupNextRequest() resets the header and context for the next request
// seqNum: next sequence number
func (c *lrpc2Client) setupNextRequest(seqNum uint32) {
	c.sequenceNum = seqNum
}

func (c *lrpc2Client) Close() error {
	return c.conn.Close()
}
