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

package stmt

import (
	"fmt"
	"sync/atomic"
)

// MySQL: COM_STMT_PREPARE Response
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html

// StatementID represents a statement ID.
type StatementID uint32

// NextStatementID returns the next statement ID.
func (s StatementID) NextStatementID() (StatementID, error) {
	if s == 0xFFFFFFFF {
		return 0, fmt.Errorf("statement ID %w", ErrOverflow)
	}

	return StatementID(atomic.AddUint32((*uint32)(&s), 1)), nil
}
