// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package binary

import (
	"bytes"
	"errors"
	"io"
)

// Reader represents a packet reader.
type Reader struct {
	io.Reader
	peekBuf []byte
}

// NewReaderWithReader returns a new packet reader with the specified reader.
func NewReaderWithReader(reader io.Reader) *Reader {
	return &Reader{
		Reader:  reader,
		peekBuf: make([]byte, 0),
	}
}

// NewReaderWithBytes returns a new packet reader with the specified byte array.
func NewReaderWithBytes(buf []byte) *Reader {
	return NewReaderWithReader(bytes.NewReader(buf))
}

// ReadByte reads a byte.
func (reader *Reader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := reader.ReadBytes(b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadBytes reads a byte array.
func (reader *Reader) ReadBytes(buf []byte) (int, error) {
	nBufSize := len(buf)
	nReadBuf := 0
	if 0 < len(reader.peekBuf) {
		nCopy := copy(buf, reader.peekBuf)
		reader.peekBuf = reader.peekBuf[nCopy:]
		nReadBuf = nCopy
	}
	if nBufSize <= nReadBuf {
		return nReadBuf, nil
	}
	nRead, err := io.ReadAtLeast(reader.Reader, buf[nReadBuf:], nBufSize-nReadBuf)
	if err != nil {
		return nReadBuf, err
	}
	nReadBuf += nRead
	return nReadBuf, err
}

// ReadNBytes reads the specified number of bytes.
func (reader *Reader) ReadNBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	nRead, err := reader.ReadBytes(buf)
	if err != nil {
		return nil, err
	}
	if nRead != n {
		return nil, newErrInvalidLength(n, nRead)
	}
	return buf, nil
}

// PeekBytes peeks a byte array.
func (reader *Reader) PeekBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	nRead, err := reader.ReadBytes(buf)
	if err != nil {
		return nil, err
	}
	if nRead != n {
		return nil, newErrInvalidLength(n, nRead)
	}
	reader.peekBuf = append(reader.peekBuf, buf...)
	return buf, nil
}

func (reader *Reader) SkipBytes(n int) error {
	for i := 0; i < n; i++ {
		_, err := reader.ReadByte()
		if err != nil {
			return err
		}
	}
	return nil
}

// PeekByte peeks a byte.
func (reader *Reader) PeekByte() (uint8, error) {
	b, err := reader.PeekBytes(1)
	if err != nil {
		return 0, err
	}
	return uint8(b[0]), nil
}

// PeekInt1 peeks a 8-bit integer.
func (reader *Reader) PeekInt1() (uint8, error) {
	return reader.PeekByte()
}

// PeekInt2 reads a 16-bit integer.
func (reader *Reader) PeekInt2() (uint16, error) {
	int16Bytes, err := reader.PeekBytes(2)
	if err != nil {
		return 0, err
	}
	return BytesToUint2(int16Bytes)
}

// PeekInt4 peeks a 32-bit integer.
func (reader *Reader) PeekInt4() (uint32, error) {
	int32Bytes, err := reader.PeekBytes(4)
	if err != nil {
		return 0, err
	}
	return BytesToUint4(int32Bytes)
}

// ReadInt1 peeks a 8-bit integer.
func (reader *Reader) ReadInt1() (uint8, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint8(b), nil
}

// ReadInt2 reads a 16-bit integer.
func (reader *Reader) ReadInt2() (uint16, error) {
	int16Bytes := make([]byte, 2)
	nRead, err := reader.ReadBytes(int16Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 2 {
		return 0, newErrInvalidLength(2, nRead)
	}
	return BytesToUint2(int16Bytes)
}

// ReadInt3 reads a 24-bit integer.
func (reader *Reader) ReadInt3() (uint32, error) {
	int24Bytes := make([]byte, 3)
	nRead, err := reader.ReadBytes(int24Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 3 {
		return 0, newErrInvalidLength(3, nRead)
	}
	return BytesToUint3(int24Bytes)
}

// ReadInt4 reads a 32-bit integer.
func (reader *Reader) ReadInt4() (uint32, error) {
	int32Bytes := make([]byte, 4)
	nRead, err := reader.ReadBytes(int32Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 4 {
		return 0, newErrInvalidLength(4, nRead)
	}
	return BytesToUint4(int32Bytes)
}

// ReadInt8 reads a 64-bit integer.
func (reader *Reader) ReadInt8() (uint64, error) {
	int64Bytes := make([]byte, 8)
	nRead, err := reader.ReadBytes(int64Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 8 {
		return 0, newErrInvalidLength(8, nRead)
	}
	return BytesToUint8(int64Bytes)
}

// ReadLengthEncodedInt reads a length encoded integer.
func (reader *Reader) ReadLengthEncodedInt() (uint64, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	switch {
	case firstByte <= 0xFB:
		return uint64(firstByte), nil
	case firstByte == 0xFC:
		v, err := reader.ReadInt2()
		return uint64(v), err
	case firstByte == 0xFD:
		v, err := reader.ReadInt3()
		return uint64(v), err
	case firstByte == 0xFE:
		v, err := reader.ReadInt8()
		return uint64(v), err
	case firstByte == NullString:
		return NullString, nil
	default:
		return 0, newErrInvalidCode("length encoded integer", uint(firstByte))
	}
}

// ReadBytesUntil reads a byte array until the specified delimiter.
func (reader *Reader) ReadBytesUntil(delim byte) ([]byte, error) {
	buf := make([]byte, 0)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == delim {
			break
		}
		buf = append(buf, b)
	}
	return buf, nil
}

// ReadNullTerminatedString reads a string until NULL.
func (reader *Reader) ReadNullTerminatedString() (string, error) {
	strBytes, err := reader.ReadBytesUntil(0x00)
	if err != nil {
		return "", err
	}
	return string(strBytes), nil
}

// ReadNullTerminatedString reads a string until NULL.
func (reader *Reader) ReadNullTerminatedBytes() ([]byte, error) {
	bytes, err := reader.ReadBytesUntil(0x00)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// ReadEOFTerminatedString reads a string until EOF.
func (reader *Reader) ReadEOFTerminatedString() (string, error) {
	buf := make([]byte, 0)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return string(buf), nil
			}
			return "", err
		}
		buf = append(buf, b)
	}
}

// ReadFixedLengthBytes reads a fixed bytes.
func (reader *Reader) ReadFixedLengthBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := reader.ReadBytes(b)
	if err != nil {
		return nil, err
	}
	if nRead != n {
		return nil, newErrInvalidLength(n, nRead)
	}
	return b, nil
}

// ReadFixedLengthString reads a fixd string.
func (reader *Reader) ReadFixedLengthString(n int) (string, error) {
	b, err := reader.ReadFixedLengthBytes(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadVariableLengthString reads a string.
func (reader *Reader) ReadVariableLengthString(n int) (string, error) {
	return reader.ReadFixedLengthString(n)
}

// ReadLengthEncodedInt reads a length encoded integer.
func (reader *Reader) ReadLengthEncodedString() (string, error) {
	n, err := reader.ReadLengthEncodedInt()
	if err != nil {
		return "", err
	}
	switch {
	case n == 0:
		return "", nil
	case n == NullString:
		return "", ErrNull
	}
	return reader.ReadFixedLengthString(int(n))
}

// ReadLengthEncodedBytes reads a length encoded bytes.
func (reader *Reader) ReadLengthEncodedBytes() ([]byte, error) {
	n, err := reader.ReadLengthEncodedInt()
	if err != nil {
		return []byte{}, err
	}
	if n == 0 {
		return []byte{}, nil
	}
	return reader.ReadFixedLengthBytes(int(n))
}
