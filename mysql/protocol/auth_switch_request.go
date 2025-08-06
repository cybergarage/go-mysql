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

import "io"

// MySQL: Protocol::AuthSwitchRequest:
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_auth_switch_request.html

// AuthSwitchRequestOption represents the AuthSwitchRequest option function.
type AuthSwitchRequestOption func(*AuthSwitchRequest)

// WithAuthSwitchRequestPluginName returns an AuthSwitchRequestOptionFn to set the plugin name.
func WithAuthSwitchRequestPluginName(pluginName string) AuthSwitchRequestOption {
	return func(req *AuthSwitchRequest) {
		req.pluginName = pluginName
	}
}

// WithAuthSwitchRequestAuthData returns an AuthSwitchRequestOptionFn to set the auth data.
func WithAuthSwitchRequestAuthData(authData string) AuthSwitchRequestOption {
	return func(req *AuthSwitchRequest) {
		req.authData = authData
	}
}

// AuthSwitchRequest represents the MySQL Protocol::AuthSwitchRequest packet.
type AuthSwitchRequest struct {
	*packet

	status     byte
	pluginName string
	authData   string
}

func newAuthSwitchRequestWithPacket(pkt *packet) *AuthSwitchRequest {
	return &AuthSwitchRequest{
		packet:     pkt,
		status:     0xFE,
		pluginName: "",
		authData:   "",
	}
}

// NewAuthSwitchRequest creates a new AuthSwitchRequest packet.
func NewAuthSwitchRequest(opts ...AuthSwitchRequestOption) *AuthSwitchRequest {
	pkt := newAuthSwitchRequestWithPacket(newPacket())
	for _, opt := range opts {
		opt(pkt)
	}
	return pkt
}

// NewAuthSwitchRequestFromReader returns a new AuthSwitchRequest from the reader.
func NewAuthSwitchRequestFromReader(reader io.Reader) (*AuthSwitchRequest, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newAuthSwitchRequestWithPacket(pktReader)

	pkt.status, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	pkt.pluginName, err = pkt.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	pkt.authData, err = pkt.ReadEOFTerminatedString()
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

// Status returns the status.
func (pkt *AuthSwitchRequest) Status() byte {
	return pkt.status
}

// PluginName returns the plugin name.
func (pkt *AuthSwitchRequest) PluginName() string {
	return pkt.pluginName
}

// AuthData returns the auth data.
func (pkt *AuthSwitchRequest) AuthData() string {
	return pkt.authData
}

// Bytes returns the packet bytes.
func (pkt *AuthSwitchRequest) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteByte(pkt.status); err != nil {
		return nil, err
	}

	if err := w.WriteNullTerminatedString(pkt.pluginName); err != nil {
		return nil, err
	}

	if err := w.WriteEOFTerminatedString(pkt.authData); err != nil {
		return nil, err
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
