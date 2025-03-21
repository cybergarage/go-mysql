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

// MySQL: Protocol Basics
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basics.html
// MySQL: MySQL Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_packets.html
// MySQL: Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_packets.html
// MariaDB protocol difference with MySQL - MariaDB Knowledge Base
// https://mariadb.com/kb/en/mariadb-protocol-difference-with-mysql/

// SequenceID represents a MySQL packet sequence ID.
type SequenceID uint8

const (
	MaxSequenceID = SequenceID(255)
)

// Next returns the next sequence ID.
func (n SequenceID) Next() SequenceID {
	return SequenceID((n + 1) % MaxSequenceID)
}
