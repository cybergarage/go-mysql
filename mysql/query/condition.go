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

package query

import (
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Condition represents a WHERE or HAVING clause.
type Condition = vitesssp.Where

// ComparisonExpr represents a two-value comparison expression.
type ComparisonExpr = vitesssp.ComparisonExpr

// AndExpr represents an AND expression.
type AndExpr = vitesssp.AndExpr

// OrExpr represents an OR expression.
type OrExpr = vitesssp.OrExpr

// XorExpr represents an XOR expression.
type XorExpr = vitesssp.XorExpr

// NotExpr represents a NOT expression.
type NotExpr = vitesssp.NotExpr

// RangeCond represents a BETWEEN or a NOT BETWEEN expression.
type RangeCond = vitesssp.RangeCond
