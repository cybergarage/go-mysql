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

// MySQL: Protocol::Query
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html

// Query represents a MySQL Query packet.
type Query struct {
	Command
	query string
}

func newQueryWithCommand(cmd Command) *Query {
	return &Query{
		Command: cmd,
		query:   "",
	}
}

// QueryOption represents a MySQL Query option.
type QueryOption func(*Query) error

// WithQuery returns a QueryOption that sets the query.
func WithQueryString(v string) QueryOption {
	return func(pkt *Query) error {
		pkt.query = v
		return nil
	}
}

// NewQuery returns a new MySQL Query packet.
func NewQuery(opts ...QueryOption) (*Query, error) {
	pkt := newQueryWithCommand(nil)
	for _, opt := range opts {
		if err := opt(pkt); err != nil {
			return nil, err
		}
	}
	return pkt, nil
}

// NewQueryFromReader returns a new MySQL Query packet from the specified reader.
func NewQueryFromReader(reader io.Reader) (*Query, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(COM_QUERY); err != nil {
		return nil, err
	}

	return NewQueryFromCommand(cmd)
}

// NewQueryFromCommand returns a new MySQL Query packet from the specified command.
func NewQueryFromCommand(cmd Command) (*Query, error) {
	var err error

	pkt := newQueryWithCommand(cmd)
	reader := cmd.Reader()

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
	w := NewWriter()

	// query
	if err := w.WriteEOFTerminatedString(pkt.query); err != nil {
		return nil, err
	}

	pkt.Command = NewCommandWith(
		COM_QUERY,
		NewPacket(
			PacketWithSequenceID(pkt.SequenceID()),
			PacketWithPayload(w.Bytes()),
		),
	)

	return pkt.Command.Bytes()
}
