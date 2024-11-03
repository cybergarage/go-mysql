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
	"io"
)

// MySQL: Generic Response Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_response_packets.html

type GenericResponse interface {
	Header() uint8
}

// GenericResponseOption represents a MySQL generic response option.
type GenericResponseOption func(*genericResponse)

// WithGenericResponseCapability returns a generic response option to set the capabilities.
func WithGenericResponseCapability(c Capability) GenericResponseOption {
	return func(res *genericResponse) {
		res.capFlags = c
	}
}

type genericResponse struct {
	capFlags Capability
}

// NewOKFromReader returns a new OK packet from the reader.
func NewGenericResponseFromReader(r io.Reader, opts ...GenericResponseOption) (GenericResponse, error) {
	reader := NewPacketReaderWith(r)
	header, err := reader.PeekByte()
	if err != nil {
		return nil, err
	}
	res := &genericResponse{
		capFlags: 0,
	}
	res.SetOptions(opts...)

	switch header {
	case okPacketHeader:
		return NewOKFromReader(reader, WithOKCapability(res.capFlags))
	case eofPacketHeader:
		return NewEOFFromReader(reader, WithEOFCapability(res.capFlags))
	case errPacketHeader:
		return NewERRFromReader(reader, WithERRCapability(res.capFlags))
	default:
	}
	return nil, newErrInvalidHeader("response", header)
}

// SetOptions sets the options.
func (res *genericResponse) SetOptions(opts ...GenericResponseOption) {
	for _, opt := range opts {
		opt(res)
	}
}
