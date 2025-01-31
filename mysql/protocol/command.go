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

// MySQL: Command Phase
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_command_phase.html

// CommandOption represents a MySQL Command option.
type CommandOption func(Command)

// WithCommandCapabilities returns a CommandOption that sets the capabilities.
func WithCommandCapability(c Capability) CommandOption {
	return func(cmd Command) {
		cmd.SetCapability(c)
	}
}

// Command represents a MySQL command.
type Command interface {
	Packet
	// Type returns the command type.
	Type() CommandType
	// IsType returns nil if the command type is the specified type.
	IsType(t CommandType) error
	// SkipPayload skips the payload.
	SkipPayload() error
	// SetCapability sets the capabilities.
	SetCapability(Capability)
	// Capability returns the capabilities.
	Capability() Capability
}

// CommandType represents a MySQL command type.
type CommandType uint8

const (
	// ComQuit: Command Quit.
	ComQuit CommandType = 0x01
	// ComInitDB: Command Init DB.
	ComInitDB CommandType = 0x02
	// ComQuery: Command Query.
	ComQuery CommandType = 0x03
	// ComFieldList: Command Field List.
	ComFieldList CommandType = 0x04
	// ComCreateDB: Command Create DB.
	ComCreateDB CommandType = 0x05
	// ComDropDB: Command Drop DB.
	ComDropDB CommandType = 0x06
	// ComRefresh: Command Refresh.
	ComRefresh CommandType = 0x07
	// ComShutdown: Command Shutdown.
	ComShutdown CommandType = 0x08
	// ComStatistics: Command Statistics.
	ComStatistics CommandType = 0x09
	// ComProcessInfo: Command Process Info.
	ComProcessInfo CommandType = 0x0a
	// ComConnect: Command Connect.
	ComConnect CommandType = 0x0b
	// ComProcessKill: Command Process Kill.
	ComProcessKill CommandType = 0x0c
	// ComDebug: Command Debug.
	ComDebug CommandType = 0x0d
	// ComPing: Command Ping.
	ComPing CommandType = 0x0e
	// ComTime: Command Time.
	ComTime CommandType = 0x0f
	// ComDelayedInsert: Command Delayed Insert.
	ComDelayedInsert CommandType = 0x10
	// ComChangeUser: Command Change User.
	ComChangeUser CommandType = 0x11
	// ComBinlogDump: Command Binlog Dump.
	ComBinlogDump CommandType = 0x12
	// ComTableDump: Command Table Dump.
	ComTableDump CommandType = 0x13
	// ComConnectOut: Command Connect Out.
	ComConnectOut CommandType = 0x14
	// ComRegisterSlave: Command Register Slave.
	ComRegisterSlave CommandType = 0x15
	// ComStmtPrepare: Command Stmt Prepare.
	ComStmtPrepare CommandType = 0x16
	// ComStmtExecute: Command Stmt Execute.
	ComStmtExecute CommandType = 0x17
	// ComStmtSendLongData: Command Stmt Send Long Data.
	ComStmtSendLongData CommandType = 0x18
	// ComStmtClose: Command Stmt Close.
	ComStmtClose CommandType = 0x19
	// ComStmtReset: Command Stmt Reset.
	ComStmtReset CommandType = 0x1a
	// ComStmeFetch: Command Stmt Fetch.
	ComStmtFetch CommandType = 0x1b
)

// String returns the string representation of the command type.
func (t CommandType) String() string {
	switch t {
	case ComQuit:
		return "ComQuit"
	case ComInitDB:
		return "ComInitDb"
	case ComQuery:
		return "ComQuery"
	case ComFieldList:
		return "ComFieldList"
	case ComCreateDB:
		return "ComCreateDb"
	case ComDropDB:
		return "ComDropDb"
	case ComRefresh:
		return "ComRefresh"
	case ComShutdown:
		return "ComShutdown"
	case ComStatistics:
		return "ComStatistics"
	case ComProcessInfo:
		return "ComProcessInfo"
	case ComConnect:
		return "ComConnect"
	case ComProcessKill:
		return "ComProcessKill"
	case ComDebug:
		return "ComDebug"
	case ComPing:
		return "ComPing"
	case ComTime:
		return "ComTime"
	case ComDelayedInsert:
		return "ComDelayedInsert"
	case ComChangeUser:
		return "ComChangeUser"
	case ComBinlogDump:
		return "ComBinlogDump"
	case ComTableDump:
		return "ComTableDump"
	case ComConnectOut:
		return "ComConnectOut"
	case ComRegisterSlave:
		return "ComRegisterSlave"
	}
	return "ComUnknown"
}

type command struct {
	cmdType CommandType
	Packet
	capFlags Capability
}

// NewCommandWith returns a new command with the specified type and packet.
func NewCommandWith(cmdType CommandType, pkt Packet, opts ...CommandOption) Command {
	cmd := &command{
		cmdType:  cmdType,
		Packet:   pkt,
		capFlags: 0,
	}
	cmd.SetOptions(opts...)
	return cmd
}

// NewCommand returns a new command.
func NewCommand(cmdType CommandType, opts ...CommandOption) Command {
	return NewCommandWith(cmdType, nil, opts...)
}

// NewCommandFromReader returns a new command from the reader.
func NewCommandFromReader(reader io.Reader, opts ...CommandOption) (Command, error) {
	var err error

	pkt, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	if pkt.PayloadLength() <= 0 {
		return nil, newErrInvalidPacketLength(pkt.PayloadLength())
	}

	// Command Type
	cmd := &command{
		cmdType:  CommandType(pkt.payload[0]),
		Packet:   pkt,
		capFlags: 0,
	}

	cmd.SetOptions(opts...)

	return cmd, nil
}

// SetOptions sets the options.
func (cmd *command) SetOptions(opts ...CommandOption) {
	for _, opt := range opts {
		opt(cmd)
	}
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

// SetCapability sets the capabilities.
func (cmd *command) SetCapability(c Capability) {
	cmd.capFlags = c
}

// Capability returns the capabilities.
func (cmd *command) Capability() Capability {
	return cmd.capFlags
}

// SkipPayload skips the payload.
func (cmd *command) SkipPayload() error {
	if cmd.Packet == nil {
		return nil
	}
	return cmd.Packet.Reader().SkipBytes(int(cmd.PayloadLength() - 1))
}
