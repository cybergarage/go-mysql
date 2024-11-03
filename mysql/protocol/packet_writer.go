// Copyright (C) 2019 The go-mysql Authors. All rights reserved.
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

// PacketWriter represents a packet writer of MySQL protocol.
type PacketWriter struct {
	*Writer
}

// NewPacketWriter returns a new packet writer.
func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		Writer: NewWriter(),
	}
}

// WriteCommandType writes a command type.
func (w *PacketWriter) WriteCommandType(cmd Command) error {
	return w.WriteByte(byte(cmd.Type()))
}

// WriteCapabilityFlags writes the capability flags.
func (w *PacketWriter) WriteCapability(c CapabilityFlag) error {
	if c.IsEnabled(ClientProtocol41) {
		return w.WriteInt4(uint32(c))
	}
	return w.WriteInt2(uint16(c >> 16))
}

// WriteFillerBytes writes the filler bytes.
func (w *PacketWriter) WriteFillerBytes(b byte, n int) error {
	for i := 0; i < n; i++ {
		if err := w.WriteByte(b); err != nil {
			return err
		}
	}
	return nil
}

// WritePacket writes a packet.
func (w *PacketWriter) WritePacket(pkt Response) error {
	pktBytes, err := pkt.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(pktBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteOK writes a OK packet.
func (w *PacketWriter) WriteOK(opts ...any) error {
	okOpts := []OKOption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			okOpts = append(okOpts, WithOKSecuenceID(v))
		case CapabilityFlag:
			okOpts = append(okOpts, WithOKCapability(v))
		}
	}
	ok, err := NewOK(okOpts...)
	if err != nil {
		return err
	}
	okBytes, err := ok.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(okBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteErr writes a ERR packet.
func (w *PacketWriter) WriteErr(opts ...any) error {
	errOpts := []ERROption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			errOpts = append(errOpts, WithERRSecuenceID(v))
		}
	}
	pkt, err := NewERR(errOpts...)
	if err != nil {
		return err
	}
	errBytes, err := pkt.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(errBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteEOF writes a EOF packet.
func (w *PacketWriter) WriteEOF(opts ...any) error {
	eofOpts := []EOFOption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			eofOpts = append(eofOpts, WithEOFCSecuenceID(v))
		case CapabilityFlag:
			eofOpts = append(eofOpts, WithEOFCapability(v))
		}
	}
	eof, err := NewEOF(eofOpts...)
	if err != nil {
		return err
	}
	eofBytes, err := eof.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(eofBytes)
	if err != nil {
		return err
	}
	return nil
}
