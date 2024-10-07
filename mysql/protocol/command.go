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

// MySQL: Command Phase
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_command_phase.html

// Command represents a MySQL command.
type Command interface {
	Packet
	// Type returns the command type.
	Type() CommandType
	// IsType returns nil if the command type is the specified type.
	IsType(t CommandType) error
	// SkipPayload skips the payload.
	SkipPayload() error
}

// CommandType represents a MySQL command type.
type CommandType uint8

const (
	// COM_QUIT: Command Quit.
	COM_QUIT CommandType = 0x01
	// COM_INIT_DB: Command Init DB.
	COM_INIT_DB CommandType = 0x02
	// COM_QUERY: Command Query.
	COM_QUERY CommandType = 0x03
	// COM_FIELD_LIST: Command Field List.
	COM_FIELD_LIST CommandType = 0x04
	// COM_CREATE_DB: Command Create DB.
	COM_CREATE_DB CommandType = 0x05
	// COM_DROP_DB: Command Drop DB.
	COM_DROP_DB CommandType = 0x06
	// COM_REFRESH: Command Refresh.
	COM_REFRESH CommandType = 0x07
	// COM_SHUTDOWN: Command Shutdown.
	COM_SHUTDOWN CommandType = 0x08
	// COM_STATISTICS: Command Statistics.
	COM_STATISTICS CommandType = 0x09
	// COM_PROCESS_INFO: Command Process Info.
	COM_PROCESS_INFO CommandType = 0x0a
	// COM_CONNECT: Command Connect.
	COM_CONNECT CommandType = 0x0b
	// COM_PROCESS_KILL: Command Process Kill.
	COM_PROCESS_KILL CommandType = 0x0c
	// COM_DEBUG: Command Debug.
	COM_DEBUG CommandType = 0x0d
	// COM_PING: Command Ping.
	COM_PING CommandType = 0x0e
	// COM_TIME: Command Time.
	COM_TIME CommandType = 0x0f
	// COM_DELAYED_INSERT: Command Delayed Insert.
	COM_DELAYED_INSERT CommandType = 0x10
	// COM_CHANGE_USER: Command Change User.
	COM_CHANGE_USER CommandType = 0x11
	// COM_BINLOG_DUMP: Command Binlog Dump.
	COM_BINLOG_DUMP CommandType = 0x12
	// COM_TABLE_DUMP: Command Table Dump.
	COM_TABLE_DUMP CommandType = 0x13
	// COM_CONNECT_OUT: Command Connect Out.
	COM_CONNECT_OUT CommandType = 0x14
	// COM_REGISTER_SLAVE: Command Register Slave.
	COM_REGISTER_SLAVE CommandType = 0x15
)

// String returns the string representation of the command type.
func (t CommandType) String() string {
	switch t {
	case COM_QUIT:
		return "COM_QUIT"
	case COM_INIT_DB:
		return "COM_INIT_DB"
	case COM_QUERY:
		return "COM_QUERY"
	case COM_FIELD_LIST:
		return "COM_FIELD_LIST"
	case COM_CREATE_DB:
		return "COM_CREATE_DB"
	case COM_DROP_DB:
		return "COM_DROP_DB"
	case COM_REFRESH:
		return "COM_REFRESH"
	case COM_SHUTDOWN:
		return "COM_SHUTDOWN"
	case COM_STATISTICS:
		return "COM_STATISTICS"
	case COM_PROCESS_INFO:
		return "COM_PROCESS_INFO"
	case COM_CONNECT:
		return "COM_CONNECT"
	case COM_PROCESS_KILL:
		return "COM_PROCESS_KILL"
	case COM_DEBUG:
		return "COM_DEBUG"
	case COM_PING:
		return "COM_PING"
	case COM_TIME:
		return "COM_TIME"
	case COM_DELAYED_INSERT:
		return "COM_DELAYED_INSERT"
	case COM_CHANGE_USER:
		return "COM_CHANGE_USER"
	case COM_BINLOG_DUMP:
		return "COM_BINLOG_DUMP"
	case COM_TABLE_DUMP:
		return "COM_TABLE_DUMP"
	case COM_CONNECT_OUT:
		return "COM_CONNECT_OUT"
	case COM_REGISTER_SLAVE:
		return "COM_REGISTER_SLAVE"
	}
	return "COM_UNKNOWN"
}

type command struct {
	cmdType CommandType
	Packet
}

// NewCommandWith returns a new command with the specified type and packet.
func NewCommandWith(cmdType CommandType, pkt Packet) Command {
	return &command{
		cmdType: cmdType,
		Packet:  pkt,
	}
}

// NewCommand returns a new command.
func NewCommand(cmdType CommandType) Command {
	return NewCommandWith(cmdType, nil)
}

// NewCommandFromReader returns a new command from the reader.
func NewCommandFromReader(reader io.Reader) (Command, error) {
	var err error

	pkt, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	// Command Type
	cmdType, err := pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	return &command{
		cmdType: CommandType(cmdType),
		Packet:  pkt,
	}, nil
}

// Type returns the command type.
func (cmd *command) Type() CommandType {
	return cmd.cmdType
}

// IsType returns nil if the command type is the specified type.
func (cmd *command) IsType(t CommandType) error {
	if cmd.cmdType != t {
		return newErrInvalidCommandType(cmd.cmdType, t)
	}
	return nil
}

// SkipPayload skips the payload.
func (cmd *command) SkipPayload() error {
	if cmd.Packet == nil {
		return nil
	}
	return cmd.Packet.Reader().SkipBytes(int(cmd.PayloadLength() - 1))
}
