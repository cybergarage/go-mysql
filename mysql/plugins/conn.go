// Copyright (C) 2019 The go-mysql Authors. All rights reserved.
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

package plugins

import (
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// Conn represents a connection.
type Conn interface {
	tracer.Context
	// Close closes the connection.
	Close() error
	// SetDatabase sets a database name.
	SetDatabase(db string)
	// Database returns a database name.
	Database() string
	// SetTimestamp sets a timestamp.
	Timestamp() time.Time
	// UUID returns a UUID.
	UUID() uuid.UUID
	// ID returns a connection ID.
	ID() uint32
	// SetSpanContext sets a span context.
	SetSpanContext(ctx tracer.Context)
	// SpanContext returns a span context.
	SpanContext() tracer.Context
}
