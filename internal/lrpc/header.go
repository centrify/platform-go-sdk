package lrpc

import (
	"bytes"
	"encoding/binary"

	"github.com/centrify/platform-go-sdk/internal/logging"
)

// type lrpc2HeaderV4 process LRPC2 V4 message header
type lrpc2HeaderV4 struct {
	magicNum    uint32
	headerLen   uint16
	version     uint32
	pid         uint64
	sequenceNum uint32
	timestamp   uint64
	msgDataLen  uint32
}

func (h *lrpc2HeaderV4) decodeHeader(buf *bytes.Buffer) error {
	var err error

	err = binary.Read(buf, lrpc2ByteOrder, &h.magicNum)
	if err != nil {
		logging.Debugf("Failed to read magic number from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.headerLen)
	if err != nil {
		logging.Debugf("Failed to read header length from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.version)
	if err != nil {
		logging.Debugf("Failed to read version from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.pid)
	if err != nil {
		logging.Debugf("Failed to read pid from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.sequenceNum)
	if err != nil {
		logging.Debugf("Failed to read sequence number from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.timestamp)
	if err != nil {
		logging.Debugf("Failed to read timestamp from header: %v", err)
		return err
	}

	err = binary.Read(buf, lrpc2ByteOrder, &h.msgDataLen)
	if err != nil {
		logging.Debugf("Failed to read message data length from header: %v", err)
		return err
	}
	return nil
}

func (h *lrpc2HeaderV4) encodeHeader(buf *bytes.Buffer) error {
	var err error

	err = binary.Write(buf, lrpc2ByteOrder, h.magicNum)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.headerLen)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.version)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.pid)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.sequenceNum)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.timestamp)
	if err != nil {
		return err
	}

	err = binary.Write(buf, lrpc2ByteOrder, h.msgDataLen)
	if err != nil {
		return err
	}

	return nil
}

func (h *lrpc2HeaderV4) verifyHeader() error {

	if h.magicNum != lrpc2MagicNum {
		return ErrLrpc2BadMagicNum
	}

	if h.headerLen != lrpc2HeaderLengthV4 {
		return ErrLrpc2BadHeaderLen
	}

	if h.version != lrpc2Version4 {
		return ErrLrpc2MsgVerMismatch
	}

	if h.msgDataLen > lrpc2MaxMsgLen {
		return ErrLrpc2BadMsgLen
	}
	return nil
}

func (h *lrpc2HeaderV4) HeaderLen() int {
	return int(lrpc2HeaderLengthV4)
}
