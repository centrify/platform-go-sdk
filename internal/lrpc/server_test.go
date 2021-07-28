package lrpc

import (
	"time"
)

var testCtxt *LrpcTestSuite

// map of messages and corresponding handlers
var msgs = map[uint16]interface{}{
	MsgEcho:  echo,
	MsgInfo:  info,
	MsgAsync: async,
	MsgSleep: doSleep,
}

func (s *LrpcTestSuite) runServer() {

	t := s.T()
	t.Log("Starting server...")
	s.wg.Add(1)
	defer s.wg.Done()

	svr, err := NewLrpc2SessionServer(testServerEndpoint, nil)
	if err != nil {
		s.setupChannel <- false
		return
	}
	t.Logf("Server created: %v", svr)
	testCtxt = s

	err = svr.RegisterMsgsByID(msgs)
	if err != nil {
		t.Fatalf("Error in registering messages: %v", err)
		s.setupChannel <- false
		return
	}

	err = svr.Start()
	if err != nil {
		t.Fatalf("Error in starting server: %v", err)
		s.setupChannel <- false
		return
	}

	s.setupChannel <- true

	// now wait for done message
	<-s.doneChannel
	err = svr.Stop()
	if err != nil {
		t.Errorf("Error in stopping server: %v", err)
		return
	}

	err = svr.Wait()
	if err != nil {
		t.Errorf("Error in waiting for server to stop: %v", err)
		return
	}

	t.Log("Exit server process")
}

func echo(ctxt SessionCtxt, args []interface{}) []interface{} {
	var ret []interface{}

	t := testCtxt.T()
	t.Log("inside function echo")

	for _, val := range args {
		ret = append(ret, val)
	}
	return ret
}

func async(ctxt SessionCtxt, args []interface{}) []interface{} {

	t := testCtxt.T()
	t.Log("inside function async")
	t.Logf("Got async messages: %v", args)
	return nil
}

func doSleep(ctxt SessionCtxt, args []interface{}) []interface{} {
	var secToSleep uint32
	var ok bool
	var ret []interface{}

	t := testCtxt.T()

	t.Log("inside function doSleep")

	if len(args) < 1 {
		// no argument specified
		secToSleep = 1
	} else {
		secToSleep, ok = args[0].(uint32)
		if !ok {
			t.Logf("Error in argument. Expect type uint32, got type %T value %v", args[0], args[0])
			secToSleep = 1
		}
	}
	sleepDuration := time.Duration(secToSleep) * time.Second
	time.Sleep(sleepDuration)
	ret = append(ret, true)
	return ret
}

func info(ctxt SessionCtxt, args []interface{}) []interface{} {
	var ret []interface{}

	errMsg := ""

	t := testCtxt.T()
	isPrivileged, err := ctxt.IsPrivileged()
	if err != nil {
		errMsg += "Cannot get caller's privilege information. "
	}

	pid, err := ctxt.GetProcessID()
	if err != nil {
		errMsg += "Cannot get process ID of caller. "
	}

	pgm, err := ctxt.GetProgram()
	if err != nil {
		errMsg += "Cannot get program of caller."
	}

	if errMsg != "" {
		// error case
		ret = append(ret, false)
		ret = append(ret, errMsg)
		t.Logf("Error encountered in getting context information: %s\n", errMsg)
	} else {
		// good case
		ret = append(ret, true)
		ret = append(ret, isPrivileged)
		ret = append(ret, pid)
		ret = append(ret, pgm)
	}
	return ret
}
