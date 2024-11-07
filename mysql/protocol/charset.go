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

// MySQL: Character Set
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_character_set.html#a_protocol_character_set

// CharSet represents a MySQL Character Set.
type CharSet uint8

const (
	// CharSetLatin1 represents the latin1 character set.
	CharSetLatin1 CharSet = 0x08
	// CharSetUTF8 represents the utf8 character set.
	CharSetUTF8 CharSet = 0x21
	// CharSetBinary represents the binary character set.
	CharSetBinary CharSet = 0x3f
)