package lrpc

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/centrify/platform-go-sdk/testutils"
	"github.com/centrify/platform-go-sdk/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LrpcTestSuite struct {
	testutils.CfyTestSuite
	endpoint     string
	setupChannel chan bool
	doneChannel  chan bool
	wg           sync.WaitGroup
}

const testServerEndpoint = "/tmp/LRPCTestEndPoint"

// test message IDs
const (
	MsgEcho         = 100
	MsgInfo         = 101
	MsgAsync        = 102
	MsgSleep        = 103
	MsgNotSupported = 104
)

func TestLrpcTestSuite(t *testing.T) {
	suite.Run(t, new(LrpcTestSuite))
}

func (s *LrpcTestSuite) SetupSuite() {

	t := s.T()
	t.Log("Setup to run server in background")
	s.setupChannel = make(chan bool)
	s.doneChannel = make(chan bool)

	go s.runServer()

	ok := <-s.setupChannel
	if !ok {
		t.Fatal("Error in setup server process, abort")
	}
	t.Log("Server setup correctly")
}

func (s *LrpcTestSuite) TearDownSuite() {
	if s.doneChannel != nil {
		s.T().Log("Shutdown server")
		s.doneChannel <- true
	}
	s.wg.Wait()
}
func (s *LrpcTestSuite) TestNoEndpoint() {
	t := s.T()
	t.Log("Testing no endpoint for lrpc")

	cl := NewLrpc2ClientSession("/etc/nosuchendpoint")
	err := cl.Connect()
	s.Assert().Error(err, "Should not be able to connect to nonexistent endpoint")
}

func (s *LrpcTestSuite) setupClient() MessageClient {
	cl := NewLrpc2ClientSession(testServerEndpoint)
	if cl == nil {
		s.T().Error("Cannot set up client connection")
	}
	err := cl.Connect()
	if err != nil {
		s.T().Errorf("Cannot connect to server: %v\n", err)
	}
	return cl
}

// TestMsgIDNotSupported tests sending unknown message ID
func (s *LrpcTestSuite) TestMsgIDNotSupported() {
	t := s.T()

	t.Log("Test using sending msgID that is not supported")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}
	req = append(req, "test")

	res, err := DoRequest(cl, MsgNotSupported, req)
	s.Assert().Error(err, "DoRequest should return error. Error: %v", err)
	s.Assert().Nil(res, "There should be no result")
}
func (s *LrpcTestSuite) TestEcho() {
	t := s.T()

	t.Log("Test echo of arguments")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}

	req = append(req, uint32(12345))
	req = append(req, true)
	req = append(req, "test string")
	req = append(req, int32(-3456))

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")

	// check size of result and individual ones
	s.Assert().Equal(len(req), len(res), "Result should be same length as request")
	if len(req) == len(res) {
		for i, v := range req {
			s.Assert().Equalf(v, res[i], "", "For element %d in result: expect: %v (type %T), got: %v (type %T)",
				i, v, v, res[i], res[i])
		}
	}
}

func (s *LrpcTestSuite) TestAsyncEcho() {
	t := s.T()

	t.Log("Test echo of arguments but don't care about results")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}

	req = append(req, "don't care")

	err := DoAsyncRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")
}

// TestKvSet tests using of key/value set as parameters and return vaules
func (s *LrpcTestSuite) TestKvSet() {
	t := s.T()

	t.Log("Test key/value sets as arguments/parameters")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}
	data := map[string]string{
		"key1":       "value1",
		"key2":       "value2",
		"null-value": "",
		"key3":       "value3",
	}

	req = append(req, data)

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")

	// check size of result and individual ones
	s.Assert().Equal(len(req), len(res), "Result should be same length as request")
	if len(req) == len(res) {
		for i, v := range req {
			same := assert.ObjectsAreEqual(v, res[i])
			s.Assert().Truef(same, "For element %d in result: expect: %v (type %T), got: %v (type %T)",
				i, v, v, res[i], res[i])
		}
	}

}

// TestStringSet tests using string set as parameters and return values
func (s *LrpcTestSuite) TestStringSet() {
	t := s.T()

	t.Log("Test using string array as arguments/parameters")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}
	data := make([]string, 10)
	data[0] = "value1"
	data[1] = "value2"
	data[2] = ""
	data[3] = "preceeded by empty string"

	req = append(req, data)

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")

	// check size of result and individual ones
	s.Assert().Equal(len(req), len(res), "Result should be same length as request")
	if len(req) == len(res) {
		for i, v := range req {
			same := assert.ObjectsAreEqual(v, res[i])
			s.Assert().Truef(same, "For element %d in result: expect: %v (type %T), got: %v (type %T)",
				i, v, v, res[i], res[i])
		}
	}
}

// TestBlob tests using blob as parameters and return values
func (s *LrpcTestSuite) TestBlob() {
	t := s.T()

	t.Log("Test using byte array as arguments/parameters")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}
	data := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = byte((i * 2) % 255)
	}

	req = append(req, data)

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")

	// check size of result and individual ones
	s.Assert().Equal(len(req), len(res), "Result should be same length as request")
	if len(req) == len(res) {
		for i, v := range req {
			same := assert.ObjectsAreEqual(v, res[i])
			s.Assert().Truef(same, "For element %d in result: expect: %v (type %T), got: %v (type %T)",
				i, v, v, res[i], res[i])
		}
	}
}

// TestNilData tests sending/receiving nil as parameter
func (s *LrpcTestSuite) TestNil() {
	t := s.T()

	t.Log("Test using nil as arguments/parameters")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}

	req = append(req, "abc")
	req = append(req, nil)
	req = append(req, uint32(123))

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().NoError(err, "DoRequest should not return error")

	// check size of result and individual ones
	s.Assert().Equal(len(req), len(res), "Result should be same length as request")
	if len(req) == len(res) {
		for i, v := range req {
			same := assert.ObjectsAreEqual(v, res[i])
			s.Assert().Truef(same, "For element %d in result: expect: %v (type %T), got: %v (type %T)",
				i, v, v, res[i], res[i])
		}
	}
}

// TestUnsupportedType tests sending unsupported types (e.g., float)
func (s *LrpcTestSuite) TestUnsupportedType() {
	t := s.T()

	t.Log("Test using unsupported type")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}

	req = append(req, uint32(123))
	req = append(req, 1.234)

	res, err := DoRequest(cl, MsgEcho, req)
	s.Assert().ErrorIs(err, ErrLrpc2TypeNotSupported)
	s.Assert().Nil(res, "No result expected")

}

// TestClientInfo verifies that server can get caller information
func (s *LrpcTestSuite) TestClientInfo() {
	t := s.T()

	t.Log("Test get caller information")

	cl := s.setupClient()
	defer cl.Close()

	res, err := DoRequest(cl, MsgInfo, nil)
	s.Assert().NoError(err, "Request should be sent with no error")

	s.Assert().GreaterOrEqualf(len(res), 2, "Result should have at least 2 elements")

	status, ok := res[0].(bool)
	s.Assert().Truef(ok, "First response must be boolean")
	s.Assert().Truef(status, "First response must be true")
	if !status {
		errmsg, _ := res[1].(string)
		s.Failf("Server return error: [%s]\n", errmsg)
	} else {
		t.Logf("Returned result: %v", res)

		isPriv, ok := res[1].(bool)
		s.Assert().Truef(ok, "Second response value must be boolean")
		runAsPriv, err := utils.RunByPrivilegedUser()
		s.Assert().NoError(err, "Should get information about whether current process is run by privileged user")
		s.Assert().Equalf(runAsPriv, isPriv, "privilege user information incorrect")

		// check process ID
		pid, ok := res[2].(int32)
		s.Assert().Truef(ok, "Returned PID should be int32")
		s.Assert().Equalf(int(pid), os.Getpid(), "Process ID should be the same")

		// check program name
		s.Assert().Equalf(len(res), 4, "Good result should have 4 elements")
		if len(res) >= 4 {
			name, ok := res[3].(string)
			s.Assert().True(ok, "Program name should be a string")
			s.Assert().Containsf(name, ".test", "Program name should contain .test")
		}
	}
}

// TestAaync verifies that async messages can be sent
func (s *LrpcTestSuite) TestAsync() {
	t := s.T()

	t.Log("Test sending async request")

	cl := s.setupClient()
	defer cl.Close()

	var req []interface{}

	req = append(req, "test message")
	req = append(req, uint32(123))

	err := DoAsyncRequest(cl, MsgAsync, req)
	s.Assert().NoError(err, "Async request should have no error")

}

// TestSingleSleep test sleep function
func (s *LrpcTestSuite) doSleep(wg *sync.WaitGroup, secToSleep uint32) {
	t := s.T()

	t.Logf("Test sending single sleep request of %v seconds\n", secToSleep)
	cl := s.setupClient()
	defer func() {
		cl.Close()
		wg.Done()
	}()

	var req []interface{}

	req = append(req, secToSleep)

	_, err := DoRequest(cl, MsgSleep, req)
	s.Assert().NoError(err, "Send request should have no error")
}

// TestConcurrent tests multiple concurrent requests
func (s *LrpcTestSuite) TestConcurrent() {

	t := s.T()

	var wg sync.WaitGroup

	startTime := time.Now()

	// send 3 different sleep requests of 1, 2, 3 seconds, this function should
	// complete a little bit over 3 seconds
	wg.Add(3)
	go s.doSleep(&wg, uint32(1))
	go s.doSleep(&wg, uint32(2))
	go s.doSleep(&wg, uint32(3))
	wg.Wait()

	elapsed := time.Since(startTime)
	t.Logf("Elapsed time: %v\n", elapsed)
	s.Assert().Greater(elapsed, 3*time.Second, "Should be longer than the longest parallel tests")
	s.Assert().Less(elapsed, 4*time.Second, "concurrent requests are not served in parallel")

}
