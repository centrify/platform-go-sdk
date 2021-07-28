package lrpc

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

// type lrpc2Server represents a LRPC2 server session.  It must implement the MessageServer interface
type lrpc2Server struct {
	lrpc2HeaderV4
	endpoint       string          // endpoint information
	listener       net.Listener    // main listener
	conn           net.Conn        // network connection
	config         *lrpc2SrvConfig // lrpc2 server config
	connected      bool            // flag to show if connected
	doneHandShake  chan bool       // channel for synchronizing handshake
	inHandShake    bool            // state to show whether handshake in progress
	handShakeError error           // handshake error

}

// lrpc2SrvConfig contains LRPC server configuration
type lrpc2SrvConfig struct {
	// LRPC2 server-side timeout when establishing connection (handshake)
	ConnectTimeout time.Duration

	// LRPC2 server-side timeout when reading data from a client request
	ReceiveTimeout time.Duration

	// LRPC2 server-side timeout when sending reply back to client
	SendTimeout time.Duration
}

func initLrpc2ServerConfig() *lrpc2SrvConfig {
	return &lrpc2SrvConfig{
		ConnectTimeout: lrpc2ConnectTimeout,
		ReceiveTimeout: lrpc2ReceiveTimeout,
		SendTimeout:    lrpc2SendTimeout,
	}
}

// initLrpc2ServerSession init a new LRPC2 server session object.
func initLrpc2ServerSession(endpoint string, config *lrpc2SrvConfig) *lrpc2Server {
	s := new(lrpc2Server)
	s.magicNum = lrpc2MagicNum
	s.headerLen = lrpc2HeaderLengthV4
	s.version = lrpc2Version4
	s.endpoint = endpoint
	s.connected = false
	s.config = config
	return s
}

// createMessageServer creates a new LRPC2 server session object.
func createMessageServer(endpoint string, acl AccessControlLevel) MessageServer {
	var err error

	config := initLrpc2ServerConfig()
	s := initLrpc2ServerSession(endpoint, config)
	if config == nil || s == nil {
		logging.Errorf("LRPC Server: Cannot init message server session for endpoint %s", endpoint)
		return nil
	}

	// start listener
	s.listener, err = StartListener(s.endpoint, acl)
	if err != nil {
		// cannot start listener
		logging.Errorf("LRPC Server: Cannot create message server session for endpoint %s: %v", endpoint, err)
		return nil
	}
	logging.Debugf("LRPC Server: New message server session [%p] created for endpoint %s", s, endpoint)
	return s
}

// Accept() accepts an incoming connection.
// Notes:
// 1. It needs to perform the LRPC handshake
// 2. It returns error if the incoming connection does not completes the handshake.  In this case, there is no new
//    socket or MessageServer object returned.  However, the caller should continue to listen to further requests for
//    the same endpoint (i.e., others should still be able to connect to it)
func (s *lrpc2Server) Accept() (MessageServer, error) {
	var err error

	if s.connected {
		return nil, ErrLrpcServerAlreadyConnected
	}

	// do accept in network layer first
	newconn, err := s.listener.Accept()
	if err != nil {
		if err == io.EOF {
			logging.Debug("Connection closed")
		}
		// cannot accept
		return nil, err
	}

	// clone a new MessageServer object for the new connection
	ret := initLrpc2ServerSession(s.endpoint, s.config)

	ret.doneHandShake = make(chan bool, 1)
	ret.inHandShake = true
	ret.conn = newconn

	// note: actual handshake is done as a goroutine....
	logging.Debugf("LRPC SERVER:  accept new connection %p. Network connection %p", ret, newconn)
	go ret.doHandShake()
	return ret, nil
}

func (s *lrpc2Server) sendBuf(data []interface{}) error {
	var err error

	buf := new(bytes.Buffer)

	// encoding data into byte buffer
	for _, v := range data {
		err = binary.Write(buf, lrpc2ByteOrder, v)
		if err != nil {
			logging.Errorf("LRPC SERVER: Failed to write to intermediate buffer: %v", err)
			return err
		}
	}
	return write(s.conn, buf.Bytes())
}

func (s *lrpc2Server) replyNACK() error {
	data := []interface{}{
		uint32(lrpc2Nack),
	}
	err := s.sendBuf(data)
	if err != nil {
		logging.Errorf("LRPC SERVER: Error in sending NACK in connection %p: %v", s, err)
	}
	return err
}

func (s *lrpc2Server) replyACK() error {
	data := []interface{}{
		uint32(lrpc2Ack),
		uint32(lrpc2MaxMsgLen),
	}
	err := s.sendBuf(data)
	if err != nil {
		logging.Errorf("LRPC SERVER: Error in sending ACK in connection %p: %v", s, err)
	}
	return err
}

func (s *lrpc2Server) setHandShakeErr(err error) {
	s.inHandShake = false
	s.connected = false
	s.handShakeError = err
	s.doneHandShake <- true // unblocker any potential caller
}

func (s *lrpc2Server) doHandShake() {
	var err error

	//
	// Set timeout
	//

	// Note that the connection timeout is for the whole handshake. So no
	// need to set again for every read/write operation.
	now := time.Now()
	err = s.conn.SetDeadline(now.Add(s.config.ConnectTimeout))

	// NOT sending NACK back if failed to set timeout. Since we cannot
	// guarantee the write operation will time out in time. So we do not
	// reply to avoid blocking indefinitely. This should be an internal
	// error that rarely occurs.
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", s.conn, err)
		return
	}

	logging.Tracef("LRPC SERVER: Start handshake for server [%p] and connection [%v] with timeout [%v]",
		s, s.conn, s.config.ConnectTimeout)

	var bytesRequest = make([]byte, lrpc2HandshakeRequestSize)
	n, err := read(s.conn, bytesRequest)
	if err != nil {
		switch t := err.(type) {
		case *net.OpError:
			if t.Err.Error() == "i/o timeout" || t.Err.Error() == "broken pipe" {
				logging.Debugf("LRPC SERVER: Connection not working. Not replying handshake for connection %p. Error: %v",
					s, err)
			} else {
				logging.Debugf("LRPC SERVER: Other connection error for connection %p: %v. Send NACK", s, err)
				s.replyNACK()
			}
		default:
			logging.Debugf("LRPC SERVER: Other non-network error for connection %p: %v.", s, err)
		}
		// update state...
		s.setHandShakeErr(err)
		return
	}

	if n != lrpc2HandshakeRequestSize {
		logging.Errorf("LRPC SERVER: Handshake request size mismatched.  Expect %d, got %d", lrpc2HandshakeRequestSize, n)
		s.replyNACK()
		// update error status
		s.setHandShakeErr(ErrMsgBadHandshakeSize)
		return
	}

	req := bytes.NewBuffer(bytesRequest)

	var protocolVersion uint32

	err = binary.Read(req, lrpc2ByteOrder, &protocolVersion)
	if err != nil {
		logging.Errorf("LRPC SERVER: Failed to read LRPC2 protocol version from handshake request for connection %p: %v", s, err)
		s.setHandShakeErr(err)
		return
	}

	if protocolVersion != lrpc2Version4 {
		logging.Errorf("LRPC SERVER: Unknown LRPC version number for connection %p: %d", s, protocolVersion)
		s.replyNACK()
		s.setHandShakeErr(ErrMsgBadVersion)
		return
	}

	logging.Tracef("LRPC SERVER: Accept handshake for connection %p.  Send ACK", s)

	err = s.replyACK()
	if err != nil {
		logging.Errorf("LRPC SERVER: Error in sending ack for connection %p: %v", s, err)
		s.setHandShakeErr(err)
		return
	}

	logging.Tracef("LRPC SERVER: connection %p is ready for requests", s)
	// everything is good...update state and return
	s.inHandShake = false
	s.connected = true
	s.doneHandShake <- true
	return
}

func (s *lrpc2Server) Close() error {
	s.connected = false
	if s.listener != nil {
		return s.listener.Close()
	}
	conn := s.conn
	if conn != nil {
		return conn.Close()
	}

	return nil
}

func (s *lrpc2Server) GetSessionCtxt() (SessionCtxt, error) {
	return NewSessionCtxt(s.conn)
}

func (s *lrpc2Server) IsNamedMessagesSupported() bool {
	return false
}

func (s *lrpc2Server) ReadRequest() (interface{}, interface{}, []interface{}, error) {
	var err error
	var hdr lrpc2HeaderV4 // for handling the received header....

	now := time.Now()
	err = s.conn.SetDeadline(now.Add(s.config.ReceiveTimeout))
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", s.conn, err)
		return nil, nil, nil, err
	}

	logging.Tracef("LRPC SERVER: ReadRequest for connection [%p] with timeout [%v]",
		s.conn, s.config.ReceiveTimeout)

	if s.inHandShake {
		// wait for handshake to be done
		if s.doneHandShake != nil {
			logging.Tracef("LRPC SERVER: Connection %p Waiting for handshake to be done first", s)
			<-s.doneHandShake
			logging.Tracef("LRPC SERVER: Connection %p Handshake is done", s)
		}
		if s.handShakeError != nil {
			logging.Debugf("LRPC SERVER: Connection %p has handshake error: %v", s, s.handShakeError)
			return nil, nil, nil, s.handShakeError
		}
	}
	if !s.connected {
		logging.Debugf("LRPC SERVER: Connection %p is not connected.", s)
		return nil, nil, nil, ErrLrpcServerNotConnected
	}

	//  good connection....

	//  read in header first
	bytesHeader := make([]byte, s.HeaderLen())
	_, err = read(s.conn, bytesHeader)

	if err != nil {
		if err == io.EOF {
			logging.Tracef("LRPC SERVER: Get EOF for connection %p.", s)
		} else {
			logging.Errorf("LRPC SERVER: Error in reading %d bytes. Connection %p: %v", s.HeaderLen(), s, err)
		}
		return nil, nil, nil, err
	}

	logging.Tracef("LRPC SERVER: Got message header for connection %p", s)
	msgHdr := bytes.NewBuffer(bytesHeader)

	err = hdr.decodeHeader(msgHdr)
	if err != nil {
		return nil, nil, nil, err
	}

	err = hdr.verifyHeader()
	if err != nil {
		logging.Debugf("LRPC SERVER: Cannot verify header for connection %p: %v", s, err)
		return nil, nil, nil, err
	}

	logging.Tracef("LRPC SERVER: Connection %p.  Header verified.  Try to read %d bytes for message body.", s, hdr.msgDataLen)

	bytesMsgData := make([]byte, hdr.msgDataLen)

	_, err = read(s.conn, bytesMsgData)

	if err != nil {
		logging.Debugf("LRPC SERVER: Connection %p.  Error in reading %d bytes of message body: %v", s, hdr.msgDataLen, err)
		return nil, nil, nil, err
	}

	msgData := bytes.NewBuffer(bytesMsgData)

	cmd, args, err := decode(msgData)
	if err != nil {
		logging.Debugf("LRPC SERVER: Connection %p.  Error in decoding incoming message: %v", s, err)
		return nil, nil, nil, err
	}

	// everything good, save the sequence number, pid and timestamp information from message header

	logging.Tracef("LRPC SERVER: Connection %p:  Read sequence number: %d", s, hdr.sequenceNum)
	return hdr, uint16(cmd), args, nil
}

func (s *lrpc2Server) WriteResponse(msg interface{}, command interface{}, results []interface{}) error {
	var err error

	now := time.Now()
	err = s.conn.SetDeadline(now.Add(s.config.SendTimeout))
	if err != nil {
		logging.Infof("Failed to set timeout for LRPC2 connection [%p]: %v", s.conn, err)
		return err
	}

	logging.Tracef("LRPC SERVER: WriteResponse to connection [%p] with timeout %v...",
		s.conn, s.config.SendTimeout)

	if !s.connected {
		logging.Debugf("LRPC SERVER: Connection %p is not connected.", s)
		return ErrLrpcServerNotConnected
	}

	hdr, ok := msg.(lrpc2HeaderV4)
	if !ok {
		logging.Errorf("LRPC SERVER:  Wrong message context for connection %p\n", s)
		return ErrMsgBadContext
	}

	cmd, ok := command.(uint16)
	if !ok {
		logging.Errorf("LRPC SERVER: wrong command type %T. Value: %v\n", command, command)
		return ErrMsgBadContext
	}

	msgData, err := encode(cmd, results)
	if err != nil {
		logging.Errorf("LRPC SERVER: Connection %p. Error in encoding results: %v", s, err)
		return err
	}

	len := uint32(msgData.Len()) + uint32(s.HeaderLen())
	bytesMsg := make([]byte, 0, len)
	msgBuf := bytes.NewBuffer(bytesMsg)

	// update the header information
	s.sequenceNum = hdr.sequenceNum
	s.pid = hdr.pid
	s.timestamp = uint64(time.Now().Unix())
	s.msgDataLen = uint32(msgData.Len())

	// encode header
	err = s.encodeHeader(msgBuf)
	if err != nil {
		logging.Errorf("LRPC SERVER: Connection %p.  Error in encoding header: %v", s, err)
		return err
	}

	// write data
	err = binary.Write(msgBuf, lrpc2ByteOrder, msgData.Bytes())
	if err != nil {
		logging.Errorf("LRPC SERVER: Connection %p.  Error in encoding message data: %v", s, err)
		return err
	}

	// send reply
	err = write(s.conn, msgBuf.Bytes())
	if err != nil {
		logging.Errorf("LRPC SERVER: Connection %p.  Error in writing reply to network: %v", s, err)
	} else {
		logging.Tracef("LRPC SERVER: Connection %p. Reply sent", s)
	}
	return err

}
