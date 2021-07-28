package lrpc

/*
DoRequest sends a command to the LRPC server and waits for the response

 Input parameters:
  cl:  Pointer to an object that implements the MessageClient interface.
  cmd: command to use
  args: arguments for command (stored as an array)

 Return values:
  results:  results as an array
  error: error information
*/
func DoRequest(cl MessageClient, cmd interface{}, args []interface{}) ([]interface{}, error) {
	err := cl.WriteRequest(cmd, args)
	if err != nil {
		return nil, err
	}

	results, err := cl.ReadResponse()
	if err != nil {
		return nil, err
	}
	return results, nil
}

/*
DoAsyncRequest sends a command to the LRPC server.  It does not wait for any response.

 Input parameters:
  cl:  Pointer to an object that implements the MessageClient interface.
  cmd: command to use
  args: arguments for command (stored as an array)

 Return value:
  error: error information
*/
func DoAsyncRequest(cl MessageClient, cmd interface{}, args []interface{}) error {
	return cl.WriteRequest(cmd, args)
}
