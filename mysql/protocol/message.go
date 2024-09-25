// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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
// MySQL: Messages
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_messages.html
// MariaDB protocol difference with MySQL - MariaDB Knowledge Base
// https://mariadb.com/kb/en/mariadb-protocol-difference-with-mysql/

// Message represents a MySQL message.
type Message struct {
	*Header
	payload []byte
}

// Payload returns the message payload.
func (msg *Message) Payload() []byte {
	return msg.payload
}
