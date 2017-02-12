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

package rtmp

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	ol "github.com/SnailTowardThesun/go-oryx-lib/logger"
	"io"
	"math"
)

const (
	NUMBER_MARKE        = 0x00
	BOOLEAN_MARKER      = 0x01
	STRING_MARKER       = 0x02
	OBJECT_MARKER       = 0x03
	MOVIECLIP_MARKER    = 0x04
	NULL_MARKER         = 0x05
	UNDEFINED_MARKER    = 0x06
	REFERENCE_MARKER    = 0x07
	ECMA_ARRAY_MARKER   = 0x08
	OBJECT_END_MARKER   = 0x09
	STRICT_ARRAY_MARKER = 0x0a
	DATE_MARKER         = 0x0b
	LONG_STRING_MARKER  = 0x0c
	UNSUPPORTED_MARKER  = 0x0d
	RECORDESET_MARKER   = 0x0e
	XML_DOCUMENT_MARKER = 0x0f
	TYPED_OBJECT_MARKER = 0x10
)

type IAMF0Item interface {
	Dumps() []byte
}

type AMF0Item struct {
	Marker  uint8
	Payload []byte
}

func (v *AMF0Item) Dumps() []byte {
	ol.T(nil, "this AMFOItem dumps")
	var buf bytes.Buffer

	buf.Write([]byte{v.Marker})
	buf.Write(v.Payload)

	return buf.Bytes()
}

func (v *AMF0Item) IsNumber() bool {
	return v.Marker == NUMBER_MARKE
}

func (v *AMF0Item) IsBoolean() bool {
	return v.Marker == BOOLEAN_MARKER
}

func (v *AMF0Item) IsString() bool {
	return v.Marker == STRING_MARKER
}

func (v *AMF0Item) IsObject() bool {
	return v.Marker == OBJECT_MARKER
}

func (v *AMF0Item) IsMovieclip() bool {
	return v.Marker == MOVIECLIP_MARKER
}

func (v *AMF0Item) IsNULL() bool {
	return v.Marker == NULL_MARKER
}

func (v *AMF0Item) IsUndefined() bool {
	return v.Marker == UNDEFINED_MARKER
}

func (v *AMF0Item) IsReference() bool {
	return v.Marker == REFERENCE_MARKER
}

func (v *AMF0Item) IsEcmaArray() bool {
	return v.Marker == ECMA_ARRAY_MARKER
}

func (v *AMF0Item) IsObjectEnd() bool {
	return v.Marker == OBJECT_END_MARKER
}

func (v *AMF0Item) IsStrictArray() bool {
	return v.Marker == STRICT_ARRAY_MARKER
}

func (v *AMF0Item) IsDate() bool {
	return v.Marker == DATE_MARKER
}

func (v *AMF0Item) IsLongString() bool {
	return v.Marker == LONG_STRING_MARKER
}

func (v *AMF0Item) IsUnSupported() bool {
	return v.Marker == UNSUPPORTED_MARKER
}

func (v *AMF0Item) IsRecordset() bool {
	return v.Marker == RECORDESET_MARKER
}

func (v *AMF0Item) IsXmlDecument() bool {
	return v.Marker == XML_DOCUMENT_MARKER
}

func (v *AMF0Item) IsTypedObject() bool {
	return v.Marker == TYPED_OBJECT_MARKER
}

type AMF0Num struct {
	AMF0Item
	Number float64
}

func NewAMF0Num(num float64) *AMF0Num {
	nu := &AMF0Num{
		Number: num,
	}

	nu.Marker = NUMBER_MARKE

	nu.Payload = make([]byte, 8)
	binary.BigEndian.PutUint64(nu.Payload, math.Float64bits(nu.Number))
	return nu
}

func ParseAMF0Num(reader io.Reader) (*AMF0Num, error) {
	nu := &AMF0Num{}
	nu.Marker = NUMBER_MARKE
	nu.Payload = make([]byte, 8)
	if n, err := io.ReadFull(reader, nu.Payload); err != nil {
		return nil, err
	} else if n != 8 {
		err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 8)
		return nil, err
	}

	nu.Number = math.Float64frombits(binary.BigEndian.Uint64(nu.Payload))
	return nu, nil
}

type AMF0Boolean struct {
	AMF0Item
	IsTrue bool
}

func NewAMF0Boolean(isTrue bool) *AMF0Boolean {
	it := &AMF0Boolean{
		IsTrue: isTrue,
	}
	it.Payload = make([]byte, 1)
	if it.IsTrue {
		it.Payload[0] = 1
	} else {
		it.Payload[0] = 0
	}

	return it
}

func ParseAMF0Boolean(reader io.Reader) (*AMF0Boolean, error) {
	it := &AMF0Boolean{}
	it.Marker = BOOLEAN_MARKER
	it.Payload = make([]byte, 1)
	if n, err := io.ReadFull(reader, it.Payload); err != nil {
		return nil, err
	} else if n != 1 {
		err = fmt.Errorf("size=%v of readed data invalid, should be %v", n, 1)
		return nil, err
	}

	if it.Payload[0] == 0 {
		it.IsTrue = false
	} else {
		it.IsTrue = true
	}

	return it, nil
}

type AMF0String struct {
	AMF0Item
	ByteLength uint16
	Bytes      []byte
}

func NewAMF0String(payload []byte) (*AMF0String, error) {
	if (len(payload)) > 65536 {
		err := fmt.Errorf("length of string in amf0 string should be less than 65535, now is %v", len(payload))
		return nil, err
	}
	it := &AMF0String{
		Bytes: payload,
	}
	it.Marker = STRING_MARKER
	it.ByteLength = uint16(len(it.Bytes))

	var buf bytes.Buffer
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, it.ByteLength)
	buf.Write(tmp)
	buf.Write(it.Bytes)

	it.Payload = buf.Bytes()

	return it, nil
}

func ParseAMF0String(reader io.Reader) (*AMF0String, error) {
	var buf bytes.Buffer
	it := &AMF0String{}

	it.Marker = STRING_MARKER

	tmp := make([]byte, 2)
	if n, err := io.ReadFull(reader, tmp); err != nil {
		return nil, err
	} else if n != 2 {
		err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 2)
		return nil, err
	}
	buf.Write(tmp)

	it.ByteLength = binary.BigEndian.Uint16(tmp)

	tmp = make([]byte, it.ByteLength)
	if n, err := io.ReadFull(reader, tmp); err != nil {
		return nil, err
	} else if n != int(it.ByteLength) {
		err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, it.ByteLength)
		return nil, err
	}

	buf.Write(tmp)
	it.Payload = buf.Bytes()
	return it, nil
}

// end of object marker
var AMF0_END_OBJECT_MARKER = []byte{0x00, 0x00, 0x09}

type AMF0Object struct {
	AMF0Item
	Properties map[string]*AMF0String
}

func (v *AMF0Object) Write(propertyKey []byte, propertyValue *AMF0String) {
	v.Properties[string(propertyKey[:])] = propertyValue
}

func (v *AMF0Object) Dumps() []byte {
	var buf bytes.Buffer
	buf.Write([]byte{v.Marker})
	for key, value := range v.Properties {
		tmp := make([]byte, 2)
		binary.BigEndian.PutUint16(tmp, uint16(len(key)))
		buf.Write(tmp)
		buf.Write([]byte(key))

		buf.Write(value.Dumps())
	}

	buf.Write(AMF0_END_OBJECT_MARKER)
	return buf.Bytes()
}

func NewAMF0Object() (*AMF0Object, error) {
	it := &AMF0Object{}
	it.Marker = OBJECT_MARKER

	return it, nil
}

func ParseAMF0Object(reader io.Reader) (*AMF0Object, error) {
	it := &AMF0Object{}
	var buf bytes.Buffer
	it.Marker = OBJECT_MARKER

	marker := make([]byte, 1)
	size := make([]byte, 2)
	for {
		if n, err := io.ReadFull(reader, size); err != nil {
			return nil, err
		} else if n != 2 {
			err = fmt.Errorf("siz=%v of readed data is invalid, should be %v", n, 2)
			return nil, err
		}
		buf.Write(size)
		if size[0] == 0x00 && size[1] == 0x00 {
			if n, err := io.ReadFull(reader, marker); err != nil {
				return nil, err
			} else if n != 1 {
				err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 1)
				return nil, err
			}
			if marker[0] == 0x09 {
				buf.Write(marker)
				break
			}
		}

		nameLength := binary.BigEndian.Uint16(size)
		name := make([]byte, nameLength)
		if n, err := io.ReadFull(reader, name); err != nil {
			return nil, err
		} else if n != int(nameLength) {
			err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, nameLength)
		}
		buf.Write(name)

		// property marker
		if n, err := io.ReadFull(reader, marker); err != nil {
			return nil, err
		} else if n != 1 {
			err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 1)
			return nil, err
		}
		buf.Write(marker)

		if marker[0] != STRING_MARKER {
			err := fmt.Errorf("property for AMF0 object should be utf-8")
			return nil, err
		}

		// property
		if pro, err := ParseAMF0String(reader); err != nil {
			return nil, err
		} else {
			it.Properties[string(name[:])] = pro
			buf.Write(pro.Dumps())
		}

	}
	it.Payload = buf.Bytes()

	return it, nil
}

type AMF0Null struct {
	AMF0Item
}

func NewAMF0Null() (*AMF0Null, error) {
	it := &AMF0Null{}
	it.Marker = NULL_MARKER

	return it, nil
}

func ParseAMF0Null(reader io.Reader) (*AMF0Null, error) {
	it := &AMF0Null{}
	it.Marker = NULL_MARKER

	return it, nil
}

type AMF0Undefined struct {
	AMF0Item
}

func NewAMF0Undefined() (*AMF0Undefined, error) {
	it := &AMF0Undefined{}
	it.Marker = UNDEFINED_MARKER

	return it, nil
}

func ParseAMF0Undefined() (*AMF0Undefined, error) {
	it := &AMF0Undefined{}
	it.Marker = UNDEFINED_MARKER

	return it, nil
}

type AMF0EcmaArray struct {
	AMF0Item
	ObjectList list.List
}

func (v *AMF0EcmaArray) Write(obj *AMF0Object) {
	v.ObjectList.PushBack(obj)
}

func (v *AMF0EcmaArray) Dumps() []byte {
	var buf bytes.Buffer

	buf.Write([]byte{v.Marker})
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(v.ObjectList.Len()))

	for i := v.ObjectList.Front(); i != nil; i = i.Next() {
		if el, ok := i.Value.(AMF0Item); ok {
			buf.Write(el.Dumps())
		}
	}

	return buf.Bytes()
}

func NewAMF0EcmaArray() (*AMF0EcmaArray, error) {
	it := &AMF0EcmaArray{}
	it.Marker = ECMA_ARRAY_MARKER

	return it, nil
}

func ParseAMF0EcmaArray(reader io.Reader) (*AMF0EcmaArray, error) {
	it := &AMF0EcmaArray{}
	it.Marker = ECMA_ARRAY_MARKER
	var buf bytes.Buffer

	size := make([]byte, 4)
	if n, err := io.ReadFull(reader, size); err != nil {
		return nil, err
	} else if n != 4 {
		err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 4)
		return nil, err
	}
	buf.Write(size)

	objCount := binary.BigEndian.Uint32(size)

	for i := uint32(0); i < objCount; i++ {
		if obj, err := ParseAMF0Object(reader); err != nil {
			return nil, err
		} else {
			it.ObjectList.PushBack(obj)
			buf.Write(obj.Dumps())
		}
	}

	it.Payload = buf.Bytes()
	return it, nil
}

type AMF0StrictArray struct {
	AMF0Item
	ArrayList list.List
}

func (v *AMF0StrictArray) Write(amf0 *AMF0Item) {
	v.ArrayList.PushBack(amf0)
}

func (v *AMF0StrictArray) Dumps() []byte {
	var buf bytes.Buffer
	buf.Write([]byte{v.Marker})

	count := make([]byte, 4)
	binary.BigEndian.PutUint32(count, uint32(v.ArrayList.Len()))

	for i := v.ArrayList.Front(); i != nil; i = i.Next() {
		if array, ok := i.Value.(AMF0Item); ok {
			buf.Write(array.Dumps())
		}
	}

	return buf.Bytes()
}

func NewAMF0StrictArray() (*AMF0StrictArray, error) {
	it := &AMF0StrictArray{}
	it.Marker = STRICT_ARRAY_MARKER

	return it, nil
}

func ParseAmf0StrictArray(reader io.Reader) (*AMF0StrictArray, error) {
	it := &AMF0StrictArray{}
	it.Marker = STRICT_ARRAY_MARKER
	var buf bytes.Buffer

	size := make([]byte, 4)
	if n, err := io.ReadFull(reader, size); err != nil {
		return nil, err
	} else if n != 4 {
		err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 4)
		return nil, err
	}
	buf.Write(size)
	listLength := binary.BigEndian.Uint32(size)

	marker := make([]byte, 1)
	for i := uint32(0); i < listLength; i++ {
		if n, err := io.ReadFull(reader, marker); err != nil {
			return nil, err
		} else if n != 1 {
			err = fmt.Errorf("size=%v of readed data is invalid, should be %v", n, 1)
			return nil, err
		}
		if marker[0] == NUMBER_MARKE {
			if el, err := ParseAMF0Num(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		} else if marker[0] == BOOLEAN_MARKER {
			if el, err := ParseAMF0Boolean(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		} else if marker[0] == STRING_MARKER {
			if el, err := ParseAMF0String(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		} else if marker[0] == NULL_MARKER {
			if el, err := ParseAMF0Null(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		} else if marker[0] == OBJECT_MARKER {
			if el, err := ParseAMF0Object(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		} else if marker[0] == ECMA_ARRAY_MARKER {
			if el, err := ParseAMF0EcmaArray(reader); err != nil {
				return nil, err
			} else {
				it.ArrayList.PushBack(el)
				buf.Write(el.Dumps())
			}
		}
	}

	it.Payload = buf.Bytes()
	return it, nil
}

type AMF0Message struct {
	ItemList list.List
}

func (v *AMF0Message) Write(it IAMF0Item) error {
	tt := &AMF0Item{}
	v.ItemList.PushBack(tt)

	ii, ok := v.ItemList.Front().Value.(IAMF0Item)
	if !ok {
		ol.E(nil, "convert failed.")
		return nil
	}
	ii.Dumps()

	return nil
}

func (v *AMF0Message) Dumps() []byte {
	var buf bytes.Buffer

	return buf.Bytes()
}

type AMF3Message struct {
}
