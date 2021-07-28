package lrpc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

// type svrBase contains common information for the message server regardless of the process model
type svrBase struct {
	// messageMapByID:	map of all handled messages to handling functions. Map by id. (note: cheat by using interface{})
	messageMapByID map[uint16]interface{}
	// connection endpoint
	connectionName string
	// base message server
	svr MessageServer
	// listenerOK?
	listenerOK bool
}

func initSvrbase(s *svrBase, name string, acl AccessControlLevel) error {
	// try to set up a message server object first
	svr := createMessageServer(name, acl)
	if svr == nil {
		return errors.New("Cannot create new message server")
	}

	s.messageMapByID = make(map[uint16]interface{})
	s.connectionName = name
	s.svr = svr
	s.listenerOK = true
	return nil
}

// RegisterMsgsByID registers handlers for multiple messages (identified by message ID)
func (s *svrBase) RegisterMsgsByID(m map[uint16]interface{}) error {
	for k, v := range m {
		s.messageMapByID[k] = v
	}
	return nil
}

// RegisterMsgByID() registers the handler for a single message by message ID
func (s *svrBase) RegisterMsgByID(k uint16, v interface{}) error {
	s.messageMapByID[k] = v
	return nil
}

// dispatch the request to the handler that is mapped
func (s *svrBase) dispatch(ctxt SessionCtxt, cmd interface{}, args []interface{}) ([]interface{}, error) {
	var f interface{}

	// note: command may come in as uint16 or uint32 but MUST fit in an uint16 value
	switch cmd.(type) {
	case uint32:
		v := cmd.(uint32)
		v1 := uint16(v)
		if v != uint32(v1) {
			// message ID is a real 32 bit number which is too big
			return nil, ErrLrpcServerCommandOutOfRange
		}
		f = s.messageMapByID[v1]
	case uint16:
		f = s.messageMapByID[cmd.(uint16)]
	default:
		return nil, fmt.Errorf("Message id type %T is not supported", cmd)
	}

	if f == nil {
		return nil, fmt.Errorf("Message [%v] not supported", cmd)
	}

	n := getFunctionName(f)
	logging.Debugf("Ready to call handler for %s[%v] connection %p", n, cmd, ctxt)
	ff := f.(func(SessionCtxt, []interface{}) []interface{})
	reply := ff(ctxt, args)
	logging.Debugf("function %s returns.", n)
	return reply, nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
