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
	"github.com/SnailTowardThesun/go-oryx-lib/rtmp"
	"math"
	"testing"
)

func TestNewAMF0Number(t *testing.T) {
	var testNumber float64
	testNumber = 64.123
	it := rtmp.NewAMF0Number(testNumber)

	if !it.IsNumber() {
		t.Error("new struct is not AMF0Number")
		return
	}

	if it.Number != testNumber {
		t.Errorf("number=%v in new amf0 number is not equal to %v", it.Number, testNumber)
		return
	}

	buf := it.Dumps()

	if buf[0] != rtmp.NUMBER_MARKER {
		t.Error("marker in dumps bytes array is invalid")
		return
	}

	dumpNumber := math.Float64frombits(binary.BigEndian.Uint64(buf[1:]))
	if dumpNumber != testNumber {
		t.Error("number in dumps bytes array is invalid")
	}

}

func TestParseAMF0Number(t *testing.T) {
	var testNumber float64
	testNumber = 64.123
	it := rtmp.NewAMF0Number(testNumber)

	// ignore the marker
	newIt, err := rtmp.ParseAMF0Number(bytes.NewReader(it.Dumps()[1:]))
	if err != nil {
		t.Error("create AMF0 Number from reader failed. err is", err)
		return
	}

	if !newIt.IsNumber() {
		t.Error("the new AMF0 item is not AMF0 Number")
		return
	}

	if newIt.Number != testNumber {
		t.Errorf("the number=%v in new AMF0 item is not equal test data=%v", newIt.Number, testNumber)
		return
	}

	if !bytes.Equal(newIt.Dumps(), it.Dumps()) {
		t.Error("the payload in new AMF0 item is not queal test data")
		return
	}
}

func TestNewAMF0Boolean(t *testing.T) {
	testBoolean := true
	it := rtmp.NewAMF0Boolean(testBoolean)

	if !it.IsBoolean() {
		t.Error("new struct is not AMF0Boolean")
		return
	}

	if !it.IsTrue {
		t.Errorf("value=%v in new amf0 boolean is not equal %v", it.IsTrue, testBoolean)
		return
	}

	buf := it.Dumps()

	if buf[0] != rtmp.BOOLEAN_MARKER {
		t.Error("Marker in dumps bytes array is invalid")
		return
	}

	if buf[1] != 0x01 {
		t.Error("value in dumps bytes array is invalid")
		return
	}
}

func TestParseAMF0Boolean(t *testing.T) {
	testBoolean := true
	it := rtmp.NewAMF0Boolean(testBoolean)

	// ignore the marker
	newItem, err := rtmp.ParseAMF0Boolean(bytes.NewReader(it.Dumps()[1:]))
	if err != nil {
		t.Error("create amf0 boolean from reader failed. err is", err)
		return
	}

	if !newItem.IsBoolean() {
		t.Error("the new amf0 item is not amf0 boolean")
		return
	}

	if !newItem.IsTrue {
		t.Errorf("the value=%v in new amf0 item is no euqal test data=%v", newItem.IsTrue, testBoolean)
		return
	}

	if !bytes.Equal(newItem.Dumps(), it.Dumps()) {
		t.Error("the dumps strings are not equal")
		return
	}
}

func TestNewAMF0String(t *testing.T) {
	testString := "test string"
	it, err := rtmp.NewAMF0String([]byte(testString))
	if err != nil {
		t.Error("create amf0 string failed")
		return
	}

	if !it.IsString() {
		t.Error("new struct is not amf0 string")
		return
	}

	if !bytes.Equal([]byte(testString), it.Bytes) {
		t.Error("string from amf0 string is invalid")
		return
	}
}

func TestParseAMF0String(t *testing.T) {
	testString := "test string"
	it, err := rtmp.NewAMF0String([]byte(testString))
	if err != nil {
		t.Error("create amf0 string failed. err is", err)
		return
	}

	newItem, err := rtmp.ParseAMF0String(bytes.NewReader(it.Dumps()[1:]))
	if err != nil {
		t.Error("parse amf0 string from reader failed. err is", err)
		return
	}

	if !newItem.IsString() {
		t.Error("the new amf0 item is not amf0 string")
		return
	}

	if !bytes.Equal([]byte(testString), newItem.Bytes) {
		t.Error("the bytes in amf0 string is not equal to test string", newItem.Bytes)
		return
	}

	if !bytes.Equal(it.Dumps(), newItem.Dumps()) {
		t.Error("the dumps string is not equal")
		return
	}
}

func TestNewAMF0Null(t *testing.T) {
	it, err := rtmp.NewAMF0Null()
	if err != nil {
		t.Error("create amf0 null failed. err is", err)
		return
	}

	if !it.IsNULL() {
		t.Error("new struct is not amf0 null")
		return
	}
}

func TestParseAMF0Null(t *testing.T) {
	it, err := rtmp.NewAMF0Null()
	if err != nil {
		t.Error("create amf0 null failed. err is", err)
		return
	}

	newItem, err := rtmp.ParseAMF0Null(bytes.NewReader(it.Dumps()))
	if err != nil {
		t.Error("parse amf0 null from reader failed. err is", err)
		return
	}

	if !newItem.IsNULL() {
		t.Error("the new amf0 struct is not null")
		return
	}

	if !bytes.Equal(it.Dumps(), newItem.Dumps()) {
		t.Error("the dumps string is not equal")
		return
	}
}
