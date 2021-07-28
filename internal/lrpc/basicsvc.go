package lrpc

import (
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

/*
Lrpc2SessionServer defines a simple session server.

Characteristics:
 1. The process model is to have a separate goroutine that serves each incoming connection.
 2. Each connection may serve one RPC message at a time.
 3. Client/Server must close the connection upon error.
 4. It embeds struct svrBase which implements basic message registration and dispatching functionalities.
*/
type Lrpc2SessionServer struct {
	svrBase
	isRunning bool
	wg        sync.WaitGroup
}

/*
NewLrpc2SessionServer creates a session server for the specific endpoint.

 Input parameters:
 endpoint:	The endpoint that the session server serves.
 acl:		The ACL data for create sesison server
*/
func NewLrpc2SessionServer(endpoint string, acl AccessControlLevel) (SessionServer, error) {

	m := new(Lrpc2SessionServer)

	// Check to see if we can create a listener first
	err := initSvrbase(&m.svrBase, endpoint, acl)
	if err != nil {
		return nil, err
	}

	m.isRunning = false

	return m, nil
}

// func handleConnection handles a connection.  It keeps on processing messages till they are done
func (s *Lrpc2SessionServer) handleConnection(msgsvr MessageServer) {

	s.wg.Add(1)
	defer s.wg.Done()

	ctxt, err := msgsvr.GetSessionCtxt()
	if err != nil {
		logging.Errorf("Cannot get peer connection context. close connection %p. Error: %v", msgsvr, err)

		// we should close this connection and return
		msgsvr.Close()
		return
	}
	logging.Debugf("Handle connection %p.  Session context: %p", msgsvr, ctxt)

	// now in the loop of receive, dispatch, reply till the incoming connection is closed

	for s.listenerOK {
		var cmd interface{}
		var args []interface{}
		var err error
		var results []interface{}
		var msgctxt interface{}

		msgctxt, cmd, args, err = msgsvr.ReadRequest()

		// todo:  break on connection closed error, otherwise continue
		if err == io.EOF {
			logging.Trace("Connection closed.")
			break
		}
		if err != nil {
			logging.Debugf("Got error in reading request: %v", err)
			break
		}

		results, err = s.dispatch(ctxt, cmd, args)
		if err != nil {
			logging.Debugf("Got error in processing request. shutdown connection: %v", err)
			// todo:  graceful error in this case
			break
		}
		if results != nil {
			err = msgsvr.WriteResponse(msgctxt, cmd, results)
			if err != nil {
				logging.Debugf("Got error in sending response; %v", err)
				break
			}
		}

	}
	logging.Debugf("Done with connection: %p", msgsvr)
	msgsvr.Close()
}

// func run() is the actual process loop for the message server
func (s *Lrpc2SessionServer) run() {
	s.wg.Add(1) // Stopping service needs to wait for this goroutine to be done too..
	defer s.wg.Done()

	for {
		conn, err := s.svr.Accept()
		if err != nil {
			// look for closed connection...this means we are done
			if strings.Contains(err.Error(), "use of closed network connection") {
				logging.Infof("Listener closed.  No more requests")
				s.isRunning = false
				return
			}
			logging.Infof("Ignore error on accept: %v", err)
		} else {
			go s.handleConnection(conn)
		}
	}
}

/*
Start starts a goroutine that starts the session server.

If the listener starts successfully, the actual process is done as a goroutine, and this function returns immediately.
Otherwise, returns error.
*/
func (s *Lrpc2SessionServer) Start() error {
	logging.Infof("Service starts on end point %s", s.connectionName)

	if !s.listenerOK {
		return errors.New("Listener not ready")
	}

	s.isRunning = true
	go s.run()
	return nil
}

/*
Stop stops the session server.  It stops receiving new connection requests.  Existing connections will complete processing of the current request and then exits.
*/
func (s *Lrpc2SessionServer) Stop() error {
	logging.Infof("Stopping service on end point %s", s.connectionName)

	err := s.svr.Close()
	if err != nil {
		logging.Errorf("Error in closing end point %s: %v", s.connectionName, err)
		s.listenerOK = false
		return err
	}
	s.listenerOK = false
	return nil
}

/*
Wait waits for all the goroutines used in the session server to exit.
*/
func (s *Lrpc2SessionServer) Wait() error {
	logging.Infof("Waiting for all goroutines to be done in end point %s", s.connectionName)
	s.wg.Wait()
	logging.Infof("All goroutines for end point %s are done.", s.connectionName)
	return nil
}
