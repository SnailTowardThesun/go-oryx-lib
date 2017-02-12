// The MIT License (MIT)
//
// Copyright (c) 2013-2016 Oryx(ossrs)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package rtmp_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/SnailTowardThesun/go-oryx-lib/rtmp"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// unit test
func TestAMF0Message_Write(t *testing.T) {
	msg := &rtmp.AMF0Message{}
	msg.Write(nil)
}

func TestNewC0C1Package(t *testing.T) {
	pkt := rtmp.NewC0C1Package()
	if pkt == nil {
		t.Error("create C0C1 package failed.")
		return
	}

	if pkt.Version != 0x03 {
		t.Error("version is error")
		return
	}

	if pkt.Timestamp == 0x00000000 {
		t.Error("timestamp is error")
		return
	}

	if pkt.Zero != 0x00000000 {
		t.Error("zero is error")
		return
	}

	if len(pkt.RandomData) < 1528 {
		t.Error("random data is error")
		return
	}

	size := len(pkt.RandomData)
	for i := 0; i < size; i++ {
		rd := pkt.RandomData[i]
		if rd > 128+32 || rd < 32 {
			t.Error("random data invalid")
			return
		}
	}
}

func TestC0C1_Dumps(t *testing.T) {
	pkt := rtmp.NewC0C1Package()

	msg := pkt.Dumps()

	if len(msg) < 1537 {
		t.Error("message from C0C1 package invalid")
		return
	}
}

func TestParseC0C1Package(t *testing.T) {
	tmp := rtmp.NewC0C1Package()

	msg := tmp.Dumps()

	pkt, err := rtmp.ParseC0C1Package(msg)
	if err != nil {
		t.Error("covert msg into c0c1 package failed. err is", err)
		return
	}

	if pkt.Version != 0x03 {
		t.Error("version is error")
		return
	}

	if pkt.Timestamp == 0x00000000 {
		t.Error("timestamp is error")
		return
	}

	if pkt.Zero != 0x00000000 {
		t.Error("zero is error")
		return
	}

	if len(pkt.RandomData) < 1528 {
		t.Error("random data is error")
		return
	}

	size := len(pkt.RandomData)
	for i := 0; i < size; i++ {
		rd := pkt.RandomData[i]
		if rd > 128+32 || rd < 32 {
			t.Error("random data invalid")
			return
		}
	}
}

func TestNewC2Package(t *testing.T) {
	c0c1 := rtmp.NewC0C1Package()
	pkt, err := rtmp.NewC2Package(c0c1.Timestamp, c0c1.RandomData)
	if err != nil {
		t.Error("create C2 package failed. err is", err)
		return
	}

	if pkt.Timestamp <= 0 {
		t.Error("timestamp in c2 package is invalid")
		return
	}

	if pkt.Timestamp2 <= 0 {
		t.Error("timestamp2 in c2 pacakge is invalid")
		return
	}

	if len(pkt.Echo) < 1528 {
		t.Error("echo in c2 package is invalid")
		return
	}

	if !bytes.Equal(pkt.Echo, c0c1.RandomData) {
		t.Error("echo in c2 should be the same as c0c1")
		return
	}
}

func TestC2_Dumps(t *testing.T) {
	c0c1 := rtmp.NewC0C1Package()
	pkt, err := rtmp.NewC2Package(c0c1.Timestamp, c0c1.RandomData)
	if err != nil {
		t.Error("create c2 package failed. err is", err)
		return
	}

	msg := pkt.Dumps()
	if len(msg) < 1536 {
		t.Error("c2 package dumps failed.")
		return
	}
}

func TestParseC2Package(t *testing.T) {
	c0c1 := rtmp.NewC0C1Package()
	originPkt, err := rtmp.NewC2Package(c0c1.Timestamp, c0c1.RandomData)
	if err != nil {
		t.Error("create c2 package failed. err is", err)
		return
	}

	msg := originPkt.Dumps()

	pkt, err := rtmp.ParseC2Package(msg)
	if err != nil {
		t.Error("parse c2 package failed. err is", err)
		return
	}

	if pkt.Timestamp <= 0 {
		t.Error("timestamp in c2 package is invalid")
		return
	}

	if pkt.Timestamp2 <= 0 {
		t.Error("timestamp2 in c2 pacakge is invalid")
		return
	}

	if len(pkt.Echo) < 1528 {
		t.Error("echo in c2 package is invalid")
		return
	}

	if !bytes.Equal(pkt.Echo, c0c1.RandomData) {
		t.Error("echo in c2 should be the same as c0c1")
		return
	}
}

func TestRtmpChunkMessage_GetCSID(t *testing.T) {
	msg := &rtmp.RtmpChunkMessage{}
	msg.BasicHeader = make([]byte, 3)

	// id = 32
	msg.BasicHeader[0] = 0x20
	if msg.GetCSID() != 32 {
		t.Error("get id failed when id=32")
		return
	}

	// id = 96
	msg.BasicHeader[0] = 0x00
	msg.BasicHeader[1] = 0x20
	if msg.GetCSID() != 96 {
		t.Error("get id failed when id=96")
		return
	}

	// id = 608
	msg.BasicHeader[0] = 0x01
	msg.BasicHeader[1] = 0x20
	msg.BasicHeader[2] = 0x02
	if msg.GetCSID() != 608 {
		t.Error("get id failed when id=544", msg.GetCSID())
		return
	}
}

func TestChunkMessage(t *testing.T) {
	chunkSize := 128
	rd := rtmp.RtmpRandomData(512 + 32)

	list, err := rtmp.ChunkMessage(rd, uint32(chunkSize), 4, 9, 1234)
	if err != nil {
		t.Error("covert byte array into chunk message failed. err is", err)
		return
	}

	if len(list) < 5 {
		t.Error(fmt.Sprintf("the number=%v of chunked messages is invalid, should be 5", len(list)))
		return
	}

	if list[0].Formt != 0 {
		t.Error(fmt.Sprintf("format=%v of first message should be 3", list[0].Formt))
		return
	}

	if len(list[0].Data) != 128 {
		t.Error(fmt.Sprintf("data length=%v in first message should be chunk size=%v", len(list[0].Data), chunkSize))
		return
	}

	if list[4].Formt != 3 {
		t.Error(fmt.Sprintf("format=%v of first message should be 3", list[0].Formt))
		return
	}

	if len(list[4].Data) != 32 {
		t.Error(fmt.Sprintf("data length=%v in last message should be %v", len(list[4].Data), 32))
		return
	}
}

func TestRtmpChunkMessage_Read(t *testing.T) {
	rd := rtmp.RtmpRandomData(512)

	list, err := rtmp.ChunkMessage(rd, 1024, 4, 9, 1234)
	if err != nil {
		t.Error("create chunk message failed. err is", err)
		return
	}

	r := bytes.NewReader(list[0].Dumps())
	msg := &rtmp.RtmpChunkMessage{}
	if err := msg.Read(r, 1024); err != nil {
		t.Error("chunk message read failed. err is", err)
		return
	}

	if !bytes.Equal(msg.BasicHeader, list[0].BasicHeader) {
		t.Error("basic header is invalid")
		return
	}

	if !bytes.Equal(msg.MessageHeader, list[0].MessageHeader) {
		t.Error("message header is invalid")
		return
	}

	if !bytes.Equal(msg.Data, list[0].Data) {
		t.Error("data is invalid")
		return
	}
}

func TestNewRtmpMsgSetChunkSize(t *testing.T) {
	chunkSize := uint32(512)
	StreamID := uint32(1024)

	pkg := rtmp.NewRtmpMsgSetChunkSize(chunkSize, StreamID)

	if pkg.MessageType != 1 {
		t.Error("message type invalid")
		return
	}

	if pkg.PayloadLength != 4 {
		t.Error("message length invalid")
		return
	}

	if pkg.Timestamp < 1 {
		t.Error("time stamp invalid")
		return
	}

	if pkg.StreamID != StreamID {
		t.Error("message id invalid")
		return
	}

	cs := binary.BigEndian.Uint32(pkg.PayLoad)
	if cs != chunkSize {
		t.Error("payload invalid")
		return
	}
}

// test case
func TestPublishStream(t *testing.T) {
}

func TestPlayStream(t *testing.T) {
}
