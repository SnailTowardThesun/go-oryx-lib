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
	"github.com/SnailTowardThesun/go-oryx-lib/rtmp"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// unit test
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

// test case
func TestPublishStream(t *testing.T) {
}

func TestPlayStream(t *testing.T) {
}
