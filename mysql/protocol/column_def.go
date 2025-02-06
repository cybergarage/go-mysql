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

// MySQL: Column Definition
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_column_definition.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/#column-definition-packet

// ColumnDef represents a MySQL Column Definition packet.
type ColumnDef interface {
	Response
	// Catalog returns the column catalog.
	Catalog() string
	// Schema returns the column schema.
	Schema() string
	// Table returns the column table.
	Table() string
	// OrgTable returns the column original table.
	OrgTable() string
	// Name returns the column name.
	Name() string
	// OrgName returns the column original name.
	OrgName() string
	// CharSet returns the column character set.
	CharSet() uint16
	// ColLength returns the column length.
	ColLength() uint32
	// ColType returns the column type.
	ColType() uint8
	// Flags returns the column flags.
	Flags() uint16
	// Decimals returns the column decimals.
	Decimals() uint8
}
