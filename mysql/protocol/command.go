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

// MySQL: Command Phase
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_command_phase.html

// CommandType represents a MySQL command type.
type CommandType uint8

const (
	// COM_QUIT: Command Quit
	Quit CommandType = 0x01
	// COM_INIT_DB: Command Init Database
	InitDB CommandType = 0x02
	// COM_QUERY: Command Query
	Query CommandType = 0x03
	// COM_FIELD_LIST: Command Field List
	FieldList CommandType = 0x04
	// COM_CREATE_DB: Command Create Database
	CreateDB CommandType = 0x05
	// COM_DROP_DB: Command Drop Database
	DropDB CommandType = 0x06
	// COM_REFRESH: Command Refresh
	Refresh CommandType = 0x07
	// COM_SHUTDOWN: Command Shutdown
	Shutdown CommandType = 0x08
	// COM_STATISTICS: Command Statistics
	Statistics CommandType = 0x09
	// COM_PROCESS_INFO: Command Process Info
	ProcessInfo CommandType = 0x0a
	// COM_CONNECT: Command Connect
	Connect CommandType = 0x0b
	// COM_PROCESS_KILL: Command Process Kill
	ProcessKill CommandType = 0x0c
	// COM_DEBUG: Command Debug
	Debug CommandType = 0x0d
	// COM_PING: Command Ping
	Ping CommandType = 0x0e
	// COM_TIME: Command Time
	Time CommandType = 0x0f
	// COM_DELAYED_INSERT: Command Delayed Insert
	DelayedInsert CommandType = 0x10
	// COM_CHANGE_USER: Command Change User
	ChangeUser CommandType = 0x11
	// COM_BINLOG_DUMP: Command Binlog Dump
	BinlogDump CommandType = 0x12
	// COM_TABLE_DUMP: Command Table Dump
	TableDump CommandType = 0x13
	// COM_CONNECT_OUT: Command Connect Out
	ConnectOut CommandType = 0x14
	// COM_REGISTER_SLAVE: Command Register Slave
	RegisterSlave CommandType = 0x15
)
