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

package query

import (
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// NewColumnDefFlagFrom converts a given Constraint into a ColumnDefFlag.
func NewColumnDefFlagFrom(c query.Constraint) (ColumnDefFlag, error) {
	cdf := ColumnDefFlag(0)
	if (c & query.PrimaryKeyConstraint) != 0 {
		cdf |= PriKeyFlag
	}

	if (c & query.NotNullConstraint) != 0 {
		cdf |= NotNullFlag
	}

	if (c & query.UniqueConstraint) != 0 {
		cdf |= UniqueKeyFlag
	}

	return cdf, nil
}
