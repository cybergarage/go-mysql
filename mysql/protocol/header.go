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

// MySQL: Messages
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_messages.html

// Header represents a MySQL message header.
type Header struct {
	length uint32
	typ    uint8
}

// Length returns the message length.
func (header *Header) Length() uint32 {
	return header.length
}

// Type returns the message type.
func (header *Header) Type() uint8 {
	return header.typ
}
