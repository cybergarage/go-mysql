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

	"github.com/cybergarage/go-safecast/safecast"
)

type Writer struct {
	bytes.Buffer
}

// NewWriter returns a new Writer.
func NewWriter() *Writer {
	return &Writer{
		Buffer: bytes.Buffer{},
	}
}

// WriteByte writes a byte.
func (w *Writer) WriteByte(b byte) error {
	return w.Buffer.WriteByte(b)
}

// WriteBytes writes a byte array.
func (w *Writer) WriteBytes(b []byte) (int, error) {
	return w.Buffer.Write(b)
}

// WriteInt1 writes a 1 byte integer.
func (w *Writer) WriteInt1(v uint8) error {
	return w.WriteByte(byte(v))
}

// WriteInt2 writes a 2 byte integer.
func (w *Writer) WriteInt2(v uint16) error {
	b := make([]byte, 2)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	_, err := w.WriteBytes(b)
	return err
}

// WriteInt3 writes a 3 byte integer.
func (w *Writer) WriteInt3(v uint32) error {
	b := make([]byte, 3)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	_, err := w.WriteBytes(b)
	return err
}

// WriteInt4 writes a 4 byte integer.
func (w *Writer) WriteInt4(v uint32) error {
	b := make([]byte, 4)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	_, err := w.WriteBytes(b)
	return err
}

// WriteInt8 writes a 64-bit integer.
func (w *Writer) WriteInt8(v uint64) error {
	b := make([]byte, 8)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	b[4] = byte((v >> 32) & 0xFF)
	b[5] = byte((v >> 40) & 0xFF)
	b[6] = byte((v >> 48) & 0xFF)
	b[7] = byte((v >> 56) & 0xFF)
	_, err := w.WriteBytes(b)
	return err
}

// WriteNullTerminatedString writes a null terminated string.
func (w *Writer) WriteNullTerminatedString(s string) error {
	_, err := w.WriteString(s)
	if err != nil {
		return err
	}
	err = w.WriteByte(0)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeFixedLengthBytes(b []byte, fb byte, n int) error {
	var wb []byte
	switch {
	case b == nil:
		wb = bytes.Repeat([]byte{fb}, n)
	case n <= len(b):
		wb = b[:n]
	default:
		wb = append(b, bytes.Repeat([]byte{fb}, n-len(b))...)
	}
	_, err := w.WriteBytes(wb)
	if err != nil {
		return err
	}
	return nil
}

// WriteFixedLengthString writes a fixed length bytes.
func (w *Writer) WriteFixedLengthBytes(b []byte, n int) error {
	return w.writeFixedLengthBytes(b, 0x00, n)
}

// WriteFixedLengthString writes a fixed length string.
func (w *Writer) WriteFixedLengthString(s string, n int) error {
	return w.writeFixedLengthBytes([]byte(s), 0x00, n)
}

// WriteVariableLengthString writes a variable length string.
func (w *Writer) WriteVariableLengthString(s string) error {
	_, err := w.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

// WriteLengthEncodedString writes a length encoded string.
func (w *Writer) WriteLengthEncodedString(s string) error {
	var n uint8
	err := safecast.ToUint8(len(s), &n)
	if err != nil {
		return err
	}
	if err := w.WriteInt1(n); err != nil {
		return err
	}
	return w.WriteFixedLengthString(s, int(n))
}

// Bytes returns the written bytes.
func (w *Writer) Bytes() []byte {
	return w.Buffer.Bytes()
}
