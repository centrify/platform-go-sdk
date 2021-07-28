package lrpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

//
// Read string data type into byte array from byte buffer.
// Return: On success, return string value or nil, and nil error
//	   On error, return nil, and ErrMsgInvalid error
func msgReadString2ByteBuffer(buf *bytes.Buffer) ([]byte, error) {
	// Read the length
	var l int32
	binary.Read(buf, lrpc2ByteOrder, &l)
	if l == -1 {
		// indicate value is zero pointer
		return nil, nil
	}

	var n int
	var err error
	sbuf := make([]byte, int(l))
	if n, err = buf.Read(sbuf); err != nil {
		return nil, err
	}

	if n != int(l) {
		return nil, ErrMsgInvalid
	}

	return sbuf, nil
}

//
// Note that LRPC2 can represent nil by setting the length to -1. So before we
// try to convert an interface{} to string in golang, we should first check if
// it is nil, e.g.:
//
// var arg0 string
//
// if args[0] != nil {
//         arg0 = args[0].(string)
// }
//
// Otherwise, when we try to convert the nil interface{} to string, golang
// panics with the following error:
// "Failed to recover: interface conversion: interface is nil, not string"
//
func decode(buf *bytes.Buffer) (uint16, []interface{}, error) {
	var err error
	var cmd uint16

	err = binary.Read(buf, lrpc2ByteOrder, &cmd)
	if err != nil {
		s := fmt.Sprintf("Failed to read LRPC2 message ID: %v", err)
		logging.Debugf(s)
		return 0, nil, errors.New(s)
	}

	ret := []interface{}{}
	for {
		// Read Msg Type first _MsgReadMsgType()
		var t byte
		err = binary.Read(buf, lrpc2ByteOrder, &t)

		if err != nil {
			s := fmt.Sprintf("Failed to read message type while processing LRPC2 message (ID: %v): %v", cmd, err)
			logging.Debugf(s)
			return 0, nil, errors.New(s)
		}

		// rtypes:
		//   1: rd_bool,
		//   2: rd_int32,
		//   3: rd_uint32,
		//   4: rd_string,
		//   5: rd_string, <-- password type, was never support from write;
		//		       MUST use rd_protected_blob for sensitive data.
		//   6: rd_blob,
		//   7: rd_string_set,
		//   8: rd_kvset,
		//   9: rd_protected_blob

		switch t {
		case msgDataTypeBool:
			var b byte
			binary.Read(buf, lrpc2ByteOrder, &b)
			ret = append(ret, b != 0)

		case msgDataTypeInt32:
			var v int32
			binary.Read(buf, lrpc2ByteOrder, &v)
			ret = append(ret, v)

		case msgDataTypeUint32:
			var v uint32
			binary.Read(buf, lrpc2ByteOrder, &v)
			ret = append(ret, v)

		case msgDataTypeString:
			var sbuf []byte
			if sbuf, err = msgReadString2ByteBuffer(buf); err != nil {
				return 0, nil, err
			}
			if sbuf != nil {
				s := string(sbuf)
				ret = append(ret, s)
			} else {
				ret = append(ret, nil)
			}

		case msgDataTypeBlob:
			var sbuf []byte
			if sbuf, err = msgReadString2ByteBuffer(buf); err != nil {
				return 0, nil, err
			}
			ret = append(ret, sbuf)

		case msgDataTypeStringSet:
			var count uint32
			binary.Read(buf, lrpc2ByteOrder, &count)
			var sset = make([]string, count)
			for i := 0; i < int(count); i++ {
				var l int32
				// Read string
				binary.Read(buf, lrpc2ByteOrder, &l)
				sbuf := buf.Next(int(l))
				s := string(sbuf)
				sset[i] = s
			}
			ret = append(ret, sset)

		case msgDataTypeKeyValueSet: // Note: this is never used at this point
			var count uint32
			binary.Read(buf, lrpc2ByteOrder, &count)
			var kvs = make(map[string]string, count)
			for i := 0; i < int(count); i++ {
				var l int32
				binary.Read(buf, lrpc2ByteOrder, &l)
				sbuf := buf.Next(int(l))
				key := string(sbuf)

				binary.Read(buf, lrpc2ByteOrder, &l)
				sbuf = buf.Next(int(l))
				value := string(sbuf)

				kvs[key] = value
			}
			ret = append(ret, kvs)

		case msgEnd:
			return cmd, ret, nil

		default:
			logging.Debugf("Internal error. Unsupported message type found while processing LRPC2 message (ID: %v)", cmd)
			return 0, nil, ErrMsgIncorrectType
		}
	}
}

func encode(cmd uint16, args []interface{}) (*bytes.Buffer, error) {
	var err error
	var buf = new(bytes.Buffer)

	err = binary.Write(buf, lrpc2ByteOrder, uint16(cmd))
	if err != nil {
		return nil, err
	}

	for _, e := range args {
		switch v := e.(type) {
		case bool:
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeBool))
			if v {
				binary.Write(buf, lrpc2ByteOrder, byte(1))
			} else {
				binary.Write(buf, lrpc2ByteOrder, byte(0))
			}

		case int32:
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeInt32))
			binary.Write(buf, lrpc2ByteOrder, v)

		case uint32:
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeUint32))
			binary.Write(buf, lrpc2ByteOrder, v)

		case []uint32:
			for _, i := range v {
				binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeUint32))
				binary.Write(buf, lrpc2ByteOrder, i)
			}

		case string: // STRING
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeString))
			binary.Write(buf, lrpc2ByteOrder, uint32(len(v)))

			if len(v) > 0 {
				buf.Write([]byte(v))
			}

		case []string: // STRING_SET
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeStringSet))
			binary.Write(buf, lrpc2ByteOrder, uint32(len(v)))
			for _, s := range v {
				binary.Write(buf, lrpc2ByteOrder, uint32(len(s)))
				buf.Write([]byte(s))
			}

		case nil: // ZERO_POINTER STRING
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeString))
			binary.Write(buf, lrpc2ByteOrder, int32(-1))

		case map[string]string: // KEY VALUE SET
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeKeyValueSet))
			binary.Write(buf, lrpc2ByteOrder, uint32(len(v)))
			for k, v := range v {
				binary.Write(buf, lrpc2ByteOrder, uint32(len(k)))
				buf.Write([]byte(k))
				binary.Write(buf, lrpc2ByteOrder, uint32(len(v)))
				buf.Write([]byte(v))
			}

		case []byte: // BLOB data type
			binary.Write(buf, lrpc2ByteOrder, byte(msgDataTypeBlob))
			binary.Write(buf, lrpc2ByteOrder, uint32(len(v)))
			if len(v) > 0 {
				buf.Write([]byte(v))
			}

		default:
			s := fmt.Sprintf("Internal error. Failed to put bytes into LRPC2 message (ID: %v, value: %v, type: %T)", cmd, v, v)
			logging.Infof(s)
			return nil, ErrLrpc2TypeNotSupported
		}
	}

	err = binary.Write(buf, lrpc2ByteOrder, byte(msgEnd))
	if err != nil {
		return nil, err
	}

	return buf, nil
}
