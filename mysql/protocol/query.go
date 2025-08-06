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
	"io"
)

// MySQL: Protocol::Query
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html

// Query represents a COM_QUERY packet.
type Query struct {
	Command

	query             string
	paramCnt          uint64
	paramSetCnt       uint64
	newParamsBindFlag uint8
	params            []*QueryParameter
	paramValues       []byte
}

// QueryParameter represents a COM_QUERY parameter.
type QueryParameter struct {
	Type uint16
	Name string
}

func newQueryWithCommand(cmd Command, opts ...QueryOption) *Query {
	q := &Query{
		Command:           cmd,
		query:             "",
		paramCnt:          0,
		paramSetCnt:       0,
		newParamsBindFlag: 0,
		params:            []*QueryParameter{},
		paramValues:       []byte{},
	}
	for _, opt := range opts {
		opt(q)
	}

	return q
}

// QueryOption represents a MySQL Query option.
type QueryOption func(*Query)

// WithQuery returns a QueryOption that sets the query.
func WithQueryString(v string) QueryOption {
	return func(pkt *Query) {
		pkt.query = v
	}
}

// WithQueryCapabilities returns a QueryOption that sets the capabilities.
func WithQueryCapability(c Capability) QueryOption {
	return func(pkt *Query) {
		pkt.SetCapability(c)
	}
}

// NewQuery returns a new MySQL Query packet.
func NewQuery(opts ...QueryOption) (*Query, error) {
	pkt := newQueryWithCommand(nil, opts...)
	return pkt, nil
}

// NewQueryFromReader returns a new MySQL Query packet from the specified reader.
func NewQueryFromReader(reader io.Reader, opts ...QueryOption) (*Query, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComQuery); err != nil {
		return nil, err
	}

	return NewQueryFromCommand(cmd, opts...)
}

// NewQueryFromCommand returns a new MySQL Query packet from the specified command.
func NewQueryFromCommand(cmd Command, opts ...QueryOption) (*Query, error) {
	var err error

	pkt := newQueryWithCommand(cmd, opts...)

	payload := cmd.Payload()
	reader := NewPacketReaderWithReader(bytes.NewBuffer(payload[1:]))

	if pkt.Capability().HasCapability(ClientQueryAttributes) {
		// parameter_count
		pkt.paramCnt, err = reader.ReadLengthEncodedInt()
		if err != nil {
			return nil, err
		}
		// parameter_set_count
		pkt.paramSetCnt, err = reader.ReadLengthEncodedInt()
		if err != nil {
			return nil, err
		}
	}

	if 0 < pkt.paramCnt {
		_, err = reader.ReadLengthEncodedBytes()
		if err != nil {
			return nil, err
		}

		pkt.newParamsBindFlag, err = reader.ReadInt1()
		if err != nil {
			return nil, err
		}

		if pkt.newParamsBindFlag == 1 {
			for range pkt.paramCnt {
				paramType, err := reader.ReadInt2()
				if err != nil {
					return nil, err
				}

				paramName, err := reader.ReadLengthEncodedString()
				if err != nil {
					return nil, err
				}

				pkt.params = append(pkt.params,
					&QueryParameter{
						Type: uint16(paramType),
						Name: paramName,
					})
			}

			pkt.paramValues, err = reader.ReadLengthEncodedBytes()
			if err != nil {
				return nil, err
			}
		}
	}

	// query
	pkt.query, err = reader.ReadEOFTerminatedString()
	if err != nil {
		return nil, err
	}

	return pkt, err
}

// Query returns the query.
func (pkt *Query) Query() string {
	return pkt.query
}

// Bytes returns the packet bytes.
func (pkt *Query) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCommandType(pkt); err != nil {
		return nil, err
	}

	if pkt.Capability().HasCapability(ClientQueryAttributes) {
		// parameter_count
		if err := w.WriteLengthEncodedInt(pkt.paramCnt); err != nil {
			return nil, err
		}
		// parameter_set_count
		if err := w.WriteLengthEncodedInt(pkt.paramSetCnt); err != nil {
			return nil, err
		}
	}

	if 0 < pkt.paramCnt {
		if err := w.WriteLengthEncodedBytes([]byte{}); err != nil {
			return nil, err
		}

		if err := w.WriteByte(pkt.newParamsBindFlag); err != nil {
			return nil, err
		}

		if pkt.newParamsBindFlag == 1 {
			for _, param := range pkt.params {
				if err := w.WriteInt2(param.Type); err != nil {
					return nil, err
				}

				if err := w.WriteLengthEncodedString(param.Name); err != nil {
					return nil, err
				}
			}

			if err := w.WriteLengthEncodedBytes(pkt.paramValues); err != nil {
				return nil, err
			}
		}
	}

	// query
	if err := w.WriteEOFTerminatedString(pkt.query); err != nil {
		return nil, err
	}

	pkt.Command = NewCommandWith(
		ComQuery,
		NewPacket(
			WithPacketSequenceID(pkt.SequenceID()),
			WithPacketPayload(w.Bytes()),
		),
	)

	return pkt.Command.Bytes()
}
