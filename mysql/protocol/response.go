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

// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html

// Response represents a response.
type Response interface {
	// SetSequenceID sets the sequence ID.
	SetSequenceID(SequenceID)
	// Bytes returns the packet bytes.
	Bytes() ([]byte, error)
}

// NewResponseWithError returns a new response with the specified error.
func NewResponseWithError(err error) (Response, error) {
	if err == nil {
		return NewOK()
	}
	return NewERR(
		WithERRMessage(err.Error()),
	)
}
