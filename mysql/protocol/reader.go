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

package protocol

import (
	"bytes"
	"errors"
	"io"

	util "github.com/cybergarage/go-mysql/mysql/encoding/bytes"
)

// Reader represents a message reader.
type Reader struct {
	io.Reader
	peekBuf []byte
}

// NewReaderWithReader returns a new message reader with the specified reader.
func NewReaderWithReader(reader io.Reader) *Reader {
	return &Reader{
		Reader:  reader,
		peekBuf: make([]byte, 0),
	}
}

// NewReaderWithBytes returns a new message reader with the specified byte array.
func NewReaderWithBytes(buf []byte) *Reader {
	return NewReaderWithReader(bytes.NewReader(buf))
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

func (reader *Reader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := reader.ReadBytes(b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (reader *Reader) PeekBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	nRead, err := reader.ReadBytes(buf)
	if err != nil {
		return nil, err
	}
	if nRead != n {
		return nil, newShortMessageError(n, nRead)
	}
	reader.peekBuf = append(reader.peekBuf, buf...)
	return buf, nil
}

// PeekInt4 reads a 32-bit integer.
func (reader *Reader) PeekInt4() (uint, error) {
	int32Bytes, err := reader.PeekBytes(4)
	if err != nil {
		return 0, err
	}
	return uint(util.BytesToUint32(int32Bytes)), nil
}

// ReadInt1 reads a 8-bit integer.
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
		return 0, newShortMessageError(2, nRead)
	}
	return util.BytesToUint16(int16Bytes), nil
}

// ReadInt3 reads a 24-bit integer.
func (reader *Reader) ReadInt3() (uint32, error) {
	int24Bytes := make([]byte, 3)
	nRead, err := reader.ReadBytes(int24Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 3 {
		return 0, newShortMessageError(3, nRead)
	}
	return util.BytesToUint24(int24Bytes), nil
}

// ReadInt4 reads a 32-bit integer.
func (reader *Reader) ReadInt4() (uint32, error) {
	int32Bytes := make([]byte, 4)
	nRead, err := reader.ReadBytes(int32Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 4 {
		return 0, newShortMessageError(4, nRead)
	}
	return util.BytesToUint32(int32Bytes), nil
}

// ReadInt8 reads a 64-bit integer.
func (reader *Reader) ReadInt8() (uint64, error) {
	int64Bytes := make([]byte, 8)
	nRead, err := reader.ReadBytes(int64Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 8 {
		return 0, newShortMessageError(8, nRead)
	}
	return util.BytesToUint64(int64Bytes), nil
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
		return nil, newShortMessageError(n, nRead)
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
