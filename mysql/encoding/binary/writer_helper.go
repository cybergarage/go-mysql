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

package binary

// WriteNullString writes a NULL string (0xFB) to the writer.
func (w *Writer) WriteNullString() error {
	return w.WriteByte(NullString)
}

// WriteTextResultsetRow writes a text resultset row.
func (w *Writer) WriteTextResultsetRowString(s *string) error {
	if s == nil {
		return w.WriteNullString()
	}

	return w.WriteLengthEncodedString(*s)
}
