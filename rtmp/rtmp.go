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

// The oryx rtmp package support bytes from/to rtmp packets.
package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	ol "github.com/SnailTowardThesun/go-oryx-lib/logger"
	"io"
	"math/rand"
	"net"
	"time"
	"math"
)

func randomData(size int) (data []byte) {
	data = make([]byte, size)

	for i := 0; i < size; i++ {
		data[i] = byte(rand.Int()%128 + 32)
	}

	return
}

type C0C1 struct {
	Version    uint8
	Zero       uint32
	Timestamp  uint32
	RandomData []byte
}

func (v *C0C1) Dumps() []byte {
	var msg bytes.Buffer
	tmp := make([]byte, 4)

	msg.WriteByte(v.Version)

	binary.BigEndian.PutUint32(tmp, v.Timestamp)
	msg.Write(tmp)
	binary.BigEndian.PutUint32(tmp, v.Zero)
	msg.Write(tmp)
	msg.Write(v.RandomData)

	return msg.Bytes()
}

func NewC0C1Package() (pkt *C0C1) {
	pkt = &C0C1{
		Version: 0x03,
		Zero:    0,
	}
	pkt.Timestamp = uint32(time.Now().Unix())
	pkt.RandomData = randomData(1528)
	return
}

func ParseC0C1Package(msg []byte) (pkt *C0C1, err error) {
	if len(msg) < 1537 {
		err = fmt.Errorf("C0C1 package size is 1537, now size=%v", len(msg))
		return nil, err
	}
	pkt = &C0C1{}

	pkt.Version = msg[0]
	if pkt.Version != 0x03 {
		err = fmt.Errorf("the version=%v of c0 is error, should be 0x03", pkt.Version)
		return nil, err
	}
	
	pkt.Timestamp = binary.BigEndian.Uint32(msg[1:5])
	pkt.Zero = binary.BigEndian.Uint32(msg[5:9])
	if pkt.Zero != 0 {
		err = fmt.Errorf("the zero=%v of c1 is error, should be 0x00000000", pkt.Zero)
		return nil, err
	}

	pkt.RandomData = msg[9:1537]

	return
}

type C2 struct {
	Timestamp  uint32
	Timestamp2 uint32
	Echo       []byte
}

func (v *C2) Dumps() []byte {
	var msg bytes.Buffer

	ts := make([]byte, 4)
	binary.BigEndian.PutUint32(ts, v.Timestamp)
	msg.Write(ts)

	binary.BigEndian.PutUint32(ts, v.Timestamp2)
	msg.Write(ts)

	msg.Write(v.Echo)

	return msg.Bytes()
}

func NewC2Package(ts uint32, rd []byte) (*C2, error) {
	if len(rd) < 1528 {
		err := fmt.Errorf("size=%v of random data in C2 invalid, should be 1528", len(rd))
		return nil, err
	}
	pkt := &C2{
		Timestamp2: ts,
		Echo:       rd,
	}
	pkt.Timestamp = uint32(time.Now().Unix())
	return pkt, nil
}

func ParseC2Package(msg []byte) (*C2, error) {
	if len(msg) < 1536 {
		err := fmt.Errorf("size=%v of c2 message is invaild, should be 1936", len(msg))
		return nil, err
	}

	pkt := &C2{}

	tmp := msg[0:4]
	pkt.Timestamp = binary.BigEndian.Uint32(tmp)

	tmp = msg[4:8]
	pkt.Timestamp2 = binary.BigEndian.Uint32(tmp)

	pkt.Echo = msg[8:1536]

	return pkt, nil
}

type S0S1 struct {
	Version    uint8
	Zero       uint32
	Timestamp  uint32
	RandomData []byte
}

func (v *S0S1) Dumps() []byte {
	var msg bytes.Buffer
	tmp := make([]byte, 4)

	msg.WriteByte(v.Version)

	binary.BigEndian.PutUint32(tmp, v.Timestamp)
	msg.Write(tmp)
	binary.BigEndian.PutUint32(tmp, v.Zero)
	msg.Write(tmp)
	msg.Write(v.RandomData)

	return msg.Bytes()
}

func NewS0S1Package() (pkt *S0S1) {
	pkt = &S0S1{
		Version:   0x03,
		Timestamp: 0,
		Zero:      0,
	}
	pkt.RandomData = randomData(1528)
	return
}

func ParseS0S1Package(msg []byte) (pkt *S0S1, err error) {
	if len(msg) < 1537 {
		err = fmt.Errorf("C0C1 package size is 1537, now size=%v", len(msg))
		return nil, err
	}
	pkt = &S0S1{}

	pkt.Version = msg[0]
	if pkt.Version != 0x03 {
		err = fmt.Errorf("the version=%v of c0 is error, should be 0x03", pkt.Version)
		return nil, err
	}
	
	pkt.Timestamp = binary.BigEndian.Uint32(msg[1:5])
	pkt.Zero = binary.BigEndian.Uint32(msg[5:9])
	pkt.RandomData = msg[9:1537]

	return
}

type S2 struct {
	Timestamp  uint32
	Timestamp2 uint32
	Echo       []byte
}

func (v *S2) Dumps() []byte {
	var msg bytes.Buffer

	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, v.Timestamp)
	msg.Write(tmp)

	binary.BigEndian.PutUint32(tmp, v.Timestamp2)
	msg.Write(tmp)

	msg.Write(v.Echo)

	return msg.Bytes()
}

func NewS2Package(ts uint32, rd []byte) (*S2, error) {
	if len(rd) < 1528 {
		err := fmt.Errorf("size=%v in s2 package is invalid, should be 1528", len(rd))
		return nil, err
	}

	pkt := &S2{
		Timestamp2: ts,
		Echo:       rd,
	}

	pkt.Timestamp2 = uint32(time.Now().Unix())
	return pkt, nil
}

func ParseS2Package(msg []byte) (*S2, error) {
	if len(msg) < 1536 {
		err := fmt.Errorf("size=%v of c2 message is invaild, should be 1936", len(msg))
		return nil, err
	}

	pkt := &S2{}

	tmp := msg[0:4]
	pkt.Timestamp = binary.BigEndian.Uint32(tmp)

	tmp = msg[4:8]
	pkt.Timestamp2 = binary.BigEndian.Uint32(tmp)

	pkt.Echo = msg[8:1536]

	return pkt, nil
}

type RtmpChunkMessage struct {
	// 1-3 bytes
	Formt           uint8
	BasicHeader     []byte
	// 0, 3, 7, 11 bytes
	MessageHeader   []byte
	// 0 or 4 byts
	ExtendTimeStamp []byte
	Data            []byte
	
	// chunk size
	ChunkSize       uint32
	
	// message info
	Timestamp       uint32
	TimestampDelta  uint32
	MessageLength   uint32
	MessageTypeId   uint8
	MessageStreamID uint32
}

func (v *RtmpChunkMessage) GetCSID()(uint32) {
	size := v.BasicHeader[0] & 0x3F
	if size > 1 {
		return uint32(size)
	}
	
	if size == 0 {
		return uint32(64 + uint32(v.BasicHeader[1]))
	}
	
	return uint32(64 + uint32(v.BasicHeader[1]) + uint32(uint32(v.BasicHeader[2]) * 256))
}

func (v *RtmpChunkMessage) Read(reader io.Reader) (error){
	// decode basic header
	buf := make([]byte, 1)
	if n, err := io.ReadFull(reader, buf);err != nil {
		ol.E(nil, "read basic header failed. err is", err)
		return err
	} else if n != 1 {
		err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
		return err
	}
	
	v.Formt = uint8(buf[0] & 0xC0)
	csId := buf[0] & 0x3F
	
	if csId == 0 {
		v.BasicHeader = make([]byte, 2)
		v.BasicHeader[0] = buf[0]
		if n, err := io.ReadFull(reader, buf);err != nil {
			ol.E(nil, "read basic header failed. err is", err)
			return err
		} else if n != 1 {
			err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
			return err
		}
		v.BasicHeader[1] = buf[0];
	} else if csId == 1 {
		v.BasicHeader = make([]byte, 3)
		v.BasicHeader[0] = buf[0]
		
		buf = make([]byte, 2)
		if n, err := io.ReadFull(reader, buf);err != nil {
			ol.E(nil, "read basic header failed. err is", err)
			return err
		} else if n != 2 {
			err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
			return err
		}
		v.BasicHeader[1], v.BasicHeader[2] = buf[0], buf[1]
	} else {
		v.BasicHeader = buf
	}
	
	messageHeaderLength := 0
	if v.Formt == 0 {
		messageHeaderLength = 11
	} else if v.Formt == 1 {
		messageHeaderLength = 7
	} else if v.Formt == 2 {
		messageHeaderLength = 3
	} else if v.Formt == 3 {
		messageHeaderLength = 0
	}
	
	v.BasicHeader = make([]byte, messageHeaderLength)
	if n, err := io.ReadFull(reader, v.MessageHeader);err != nil {
		ol.E(nil, "read basic header failed. err is", err)
		return err
	} else if n != messageHeaderLength {
		err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
		return err
	}
	
	if v.Formt == 0 {
		v.Timestamp = binary.BigEndian.Uint32(v.MessageHeader[0:3])
		v.MessageLength = binary.BigEndian.Uint32(v.MessageHeader[3:6])
		v.MessageTypeId = uint8(v.MessageHeader[6])
		v.MessageStreamID = binary.BigEndian.Uint32(v.MessageHeader[7:11])
	} else if v.Formt == 1 {
		v.TimestampDelta = binary.BigEndian.Uint32(v.MessageHeader[0:3])
		v.MessageLength = binary.BigEndian.Uint32(v.MessageHeader[3:6])
		v.MessageTypeId = uint8(v.MessageHeader[6])
	} else if v.Formt == 2 {
		v.TimestampDelta = binary.BigEndian.Uint32(v.MessageHeader[0:3])
	}
	
	msgLength := v.MessageLength
	if v.MessageLength > v.ChunkSize {
		msgLength = v.ChunkSize
	}
	
	v.Data = make([]byte, msgLength)
	if n, err := io.ReadFull(reader, v.Data);err != nil {
		ol.E(nil, "read basic header failed. err is", err)
		return err
	} else if n != msgLength {
		err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
		return err
	}
	
	return nil
}

func ChunkMessage(msg []byte, chunkSize uint32, csId uint32, msgType uint8, streamId uint32) (list []RtmpChunkMessage, err error) {
	num := math.Ceil(len(msg) / chunkSize)
	
	buf := bytes.NewBuffer(msg)
	
	list = make([]RtmpChunkMessage, num)
	list[0].Formt = 0
	list[0].Timestamp = uint32(time.Now().Unix())
	list[0].MessageLength = uint32(len(msg))
	list[0].MessageTypeId = msgType
	list[0].MessageStreamID = streamId
	
	if num > 1 {
		list[0].Data = make([]byte, chunkSize)
		buf.Write(list[0].Data)
		for i :=1; i < num; i++ {
			list[i].Formt = 3
			if (buf.Len() > chunkSize) {
				list[i].Data = make([]byte, chunkSize)
				buf.Write(list[i].Data)
			} else {
				list[i].Data = make([]byte, buf.Len())
				buf.Write(list[i].Data)
			}
		}
	} else {
		list[0].Data = make([]byte, len(msg))
		buf.Write(list[0].Data)
	}
	
	return
}

type RtmpClient interface {
	// initialize the client
	initialize(url string) error
	// handshake
	handshake() error
	// connect to server
	connect() error
	// play stream
	Play() error
	// publish stream
	Publish() error
	// send package to server
	Send() error
	// receive package from server
	Recv() error
	// close the connection
	Close() error
}

type SimpleRtmpClient struct {
	conn net.Conn
}

func NewSimpleRtmpClient(u string) (RtmpClient, error) {
	v := &SimpleRtmpClient{}

	if err := v.initialize(u); err != nil {
		ol.E(nil, "initialize the rtmp client failed. err is", err)
		return nil, err
	}

	if err := v.handshake(); err != nil {
		ol.E(nil, "do handshake with server failed. err is", err)
		return nil, err
	}

	if err := v.connect(); err != nil {
		ol.E(nil, "connect to server failed. err is", err)
		return nil, err
	}
	return v, nil
}

func (v *SimpleRtmpClient) url_parse(u string) error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) initialize(u string) error {
	var err error

	if err = v.url_parse(u); err != nil {
		ol.E(nil, "parse url failed. err is", err)
		return err
	}
	
	v.conn, err = net.Dial("tcp", "127.0.0.1:1935")
	if err != nil {
		ol.E(nil, "connect to server failed. err is", err)
		return err
	}

	return nil
}

func (v *SimpleRtmpClient) handshake() error {
	// send c0c1
	c0c1Pkg := NewC0C1Package()
	if nn, err := v.conn.Write(c0c1Pkg.Dumps()); err != nil {
		ol.E(nil, "send c0c1 failed. err is", err)
		return err
	} else if nn != len(c0c1Pkg.Dumps()) {
		err = fmt.Errorf("send c0c1 failed, size=%v of sended size is not equal %v", nn, len(c0c1Pkg.Dumps()))
		return err
	}

	// recv s0s1
	s0s1Msg := make([]byte, 1537)

	if nn, err := io.ReadFull(v.conn, s0s1Msg); err != nil {
		ol.E(nil, "read s0s1 failed. err is", err)
		return err
	} else if nn != 1537 {
		err = fmt.Errorf("size=%v of s0s1 is invalid, should be 1537", nn)
		return err
	}
	
	s0s1Pkg, err := ParseS0S1Package(s0s1Msg)
	if err != nil {
		ol.E(nil, "parse s0s1 message failed. err is", err)
		return err
	}
	
	// send c2
	c2, err := NewC2Package(s0s1Pkg.Timestamp, s0s1Pkg.RandomData)
	if err != nil {
		ol.E(nil, "create c2 failed. err is", err)
		return err
	}
	
	if nn, err := v.conn.Write(c2.Dumps()); err != nil {
		ol.E(nil, "send c2 failed. err is", err)
		return err
	} else if nn != len(c2.Dumps()) {
		err = fmt.Errorf("send c2 failed. size=%v of sended size is no equal %v", nn, len(c2.Dumps()))
	}
	
	// recv s2
	s2Msg := make([]byte, 1536)
	
	if nn, err := io.ReadFull(v.conn, s2Msg); err != nil {
		ol.E(nil, "read S2 failed. err is", err)
		return err
	} else if nn != 1536 {
		err = fmt.Errorf("size=%v of S2 is invalid, should be 1536", nn)
		return err
	}

	s2Pkg, err := ParseS2Package(s2Msg)
	if err != nil {
		ol.E(nil, "parse S2 package failed. err is", err)
		return err
	}
	
	if !bytes.Equal(c0c1Pkg.RandomData, s2Pkg.Echo) {
		err := fmt.Errorf("random in s2 is not euqal to that in c0c1")
		return err
	}
	
	ol.T(nil, "do simple handshake successfully")
	return nil
}

func (v *SimpleRtmpClient) connect() error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) Play() error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) Publish() error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) Send() error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) Recv() error {
	// TODO:FIXME: implement this function
	return nil
}

func (v *SimpleRtmpClient) Close() error {
	// TODO:FIXME: implement this function
	return nil
}

// TODO:FIXME: implement complex rtmp client
