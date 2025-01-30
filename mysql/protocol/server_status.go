// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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

// MySQL: Server Status Flag
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysql__com_8h.html

// ServerStatus represents a MySQL server status flag.
type ServerStatus uint16

const (
	ServerStatusInTrans            ServerStatus = 1
	ServerStatusAutocommit         ServerStatus = 2
	ServerMoreResultsExists        ServerStatus = 8
	ServerQueryNoGoodIndexUsed     ServerStatus = 16
	ServerQueryNoIndexUsed         ServerStatus = 32
	ServerStatusCursorExists       ServerStatus = 64
	ServerStatusLastRowSent        ServerStatus = 128
	ServerStatusDBDropped          ServerStatus = 256
	ServerStatusNoBackslashEscapes ServerStatus = 512
	ServerStatusMetadataChanged    ServerStatus = 1024
	ServerQueryWasSlow             ServerStatus = 2048
	ServerPsOutParams              ServerStatus = 4096
	ServerStatusInTransReadonly    ServerStatus = 8192
	ServerSessionStateChanged      ServerStatus = (1 << 14)
)

// IsEnabled returns true if the status flag is enabled.
func (statFlag ServerStatus) IsEnabled(flag ServerStatus) bool {
	return (statFlag & flag) != 0
}
